let conn;
const gameData = {
  username: '',
  currentRoom: '',
  life: 20,
};
window.onload = function () {
  document.getElementById('login-form').onsubmit = onSubmitLogin;
  document.getElementById('room-selection').onsubmit = onSubmitRoomSelection;
  document.getElementById('opponent-select').onsubmit = onSubmitShoot;
};

/** Ingress/egress event types **/
const eventTypes = {
  // Out
  changeRoom: 'CHANGE_ROOM',
  shoot: 'SHOOT',
  // In
  receiveImpact: 'RECEIVE_IMPACT',
  notifyImpact: 'NOTIFY_IMPACT',
  failedImpact: 'FAILED_IMPACT',
  roomList: 'ROOM_LIST',
  roomUsers: 'ROOM_USERS',
};
/**
 * @typedef {Object} WsEvent
 * @property {string} type
 * @property {string} payload optional
 */

/**
 * Event is a eneral wrapper to Recenve / Semd  messages
 * @property {string} type
 * @property {string} payload
 */
class WsEvent {
  /**
   * @param {string} type
   * @param payload
   */
  constructor(type, payload) {
    this.type = type;
    this.payload = payload;
  }

  toString() {
    return JSON.stringify({ type: this.type, payload: this.payload });
  }
}

/** The shoot button was clicked **/
class SendShootEvent {
  /**
   * @param {string} target
   **/
  constructor(target) {
    this.target = target;
  }
}

class ReceiveImpactEvent {
  /**
   * @param {string} attacker
   * @param {string} target
   * @param {int} newLife
   **/
  constructor(attacker, target, newLife) {
    this.attacker = attacker;
    this.target = target;
    this.newLife = newLife;
  }
}

class ChangeRoomEvent {
  /**
   * @param {string} room
   **/
  constructor(room) {
    this.room = room;
  }
}

class RoomListEvent {
  /**
   * @param {string[]} rooms
   **/
  constructor(rooms) {
    this.rooms = rooms;
  }
}

class RoomUsersEvent {
  /**
   * @param {string} roomm
   * @param {string[]} users
   **/
  constructor(room, users) {
    this.room = room;
    this.users = users;
  }
}

function onSubmitLogin() {
  const username = document.getElementById('username').value;
  if (!username) {
    appendToGameLogs('username cannot be empty.');
    return;
  }
  connectWebSocket(username);
  return false;
}

function onSubmitRoomSelection() {
  const newRoom = document.getElementById('room_name').value;
  changeRoom(newRoom);
  return false;
}

function onSubmitShoot() {
  const target = document.getElementById('opponent-target')?.value;
  sendShootTarget(target);
  return false;
}

function connectWebSocket(username) {
  const websocketSupported = window['WebSocket'];
  if (websocketSupported) {
    conn = new WebSocket(
      `ws://${document.location.host}/ws?username=${username}`
    );
    conn.onopen = () => changeConnectionTitle('Connected');
    conn.onclose = () => changeConnectionTitle('Disconnected');
    conn.onmessage = (e) => e.data && handleWsMessage(e.data);
  } else {
    alert('Websockets not supported');
  }
}

/**
 * @param {string} value)
 */
function changeConnectionTitle(value) {
  document.getElementById('conn-status').innerHTML = value;
}

/**
 * @param {string} message
 */
function appendToGameLogs(message) {
  document.getElementById('events-log').innerHTML +=
    '\n<div>' + message + '</div>';
}

function updateGameLife() {
  document.getElementById('game-life').innerHTML = gameData.life;
}
/**
 * @param {string} name
 */
function updateGameRoom(name) {
  document.getElementById('room-name').innerHTML = name;
}

/**
 * @param {string} newRoom
 **/
function changeRoom(newRoom) {
  const currentRoom = document.getElementById('room-name');
  if (newRoom === '' || currentRoom == newRoom) {
    return;
  }
  const payload = new ChangeRoomEvent(newRoom);
  sendWsEvent(eventTypes.changeRoom, payload);
}

/**
 * @param {WsEvent} data
 **/
function routeEvent({ type, payload }) {
  if (type === eventTypes.receiveImpact) {
    handleReceiveInpact(payload);
    return;
  } else if (type === eventTypes.notifyImpact) {
    handleNotifyInpact(payload);
    return;
  } else if (type === eventTypes.failedImpact) {
    handlerErrorMessage(payload);
  } else if (type === eventTypes.roomList) {
    handleUpdateRooms(payload);
  } else if (type === eventTypes.roomUsers) {
    handleUpdateAvailableTargets(payload);
  } else {
    appendToGameLogs('Received unhandled message: ' + type);
  }
}

/**
 * @param {ReceiveImpactEvent} payload
 **/
function handleReceiveInpact(payload) {
  const impactData = Object.assign(new ReceiveImpactEvent(), payload);
  const logMessage = `Received impact from: ${impactData.attacker}`;
  gameData.life = impactData.newLife;
  appendToGameLogs(logMessage);
  updateGameLife();
}
/**
 * @param {ReceiveImpactEvent} payload
 **/
function handleNotifyInpact(payload) {
  const impactData = Object.assign(new ReceiveImpactEvent(), payload);
  let logMessage = `Impact from: ${impactData.attacker} to: ${impactData.target};`;
  if (impactData.newLife === 0) {
    logMessage += `<br />${impactData.target} is down.`;
  }
  appendToGameLogs(logMessage);
}

/**
 * @param {RoomListEvent} payload
 **/
function handleUpdateRooms(payload) {
  console.log(payload);
  const roomOptions = payload.rooms
    .map((room) => `<option value="${room}">${room}</option>`)
    .join('');
  console.log(payload.rooms, roomOptions);
  document.getElementById('room_name').innerHTML = roomOptions;
}
/**
 * @param {RoomUsersEvent} payload
 **/
function handleUpdateAvailableTargets(payload) {
  if (payload.room != gameData.currentRoom) {
    updateGameRoom(payload.room);
    gameData.currentRoom = payload.room;
  }
  const userOptions = payload.users
    .map((username) => `<option value="${username}">${username}</option>`)
    .join('');

  document.getElementById('opponent-target').innerHTML = userOptions;
}

function handlerErrorMessage(payload) {
  appendToGameLogs(payload.error);
}

/**
 * @param {string} target
 */
function sendShootTarget(target) {
  if (!target) {
    appendToGameLogs('Shoot: Selected none.');
    return;
  } else if (target == 'any') {
    appendToGameLogs('Shoot: <b>Any</b> is not yet available.');
    return;
  }

  const payload = new SendShootEvent(target);
  sendWsEvent(eventTypes.shoot, payload);
}

/**
 * @param {string} data
 */
function handleWsMessage(data) {
  const eventData = JSON.parse(data);
  const event = Object.assign(new WsEvent(), eventData);
  routeEvent(event);
}

/**
 * @param {string} eventType
 * @param {Object} payload
 */
function sendWsEvent(eventType, payload) {
  const message = new WsEvent(eventType, payload);
  conn.send(message.toString());
}
