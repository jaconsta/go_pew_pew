let conn;
window.onload = function() {
  document.getElementById("login-form").onsubmit = onSubmitLogin;
  document.getElementById("room-selection").onsubmit = onSubmitRoomSelection;
  document.getElementById("opponent-select").onsubmit = onSubmitShoot;
}

/** Ingress/egress event types **/
const eventTypes = {
  changeRoom: "CHANGE_ROOM",
  shoot: "SHOOT",

  receiveImpact: "RECEIVE_IMPACT",
}
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
    return JSON.stringify({ type: this.type, payload: this.payload })
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
   **/
  constructor(attacker) {
    this.attacker = attacker;
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

function onSubmitLogin() {
  const username = document.getElementById("username").value;
  if (!username) {
    appendToGameLogs("username cannot be empty.")
    return
  }
  connectWebSocket(username);
  return false;
}

function onSubmitRoomSelection() {
  const newRoom = document.getElementById("room_name").value;
  changeRoom(newRoom);
  return false;
}

function onSubmitShoot() {
  const target = document.getElementById("opponent-target")?.value;
  sendShootTarget(target);
  return false;
}

function connectWebSocket(username) {
  const websocketSupported = window["WebSocket"];
  if (websocketSupported) {
    conn = new WebSocket(`ws://${document.location.host}/ws?username=${username}`);
    conn.onopen = () => changeConnectionTitle("Connected");
    conn.onclose = () => changeConnectionTitle("Disconnected");
    conn.onmessage = (e) => e.data && handleWsMessage(e.data);
  } else {
    alert("Websockets not supported");
  }
}

/**
  * @param {string} value)
*/
function changeConnectionTitle(value) {
  document.getElementById("conn-status").innerHTML = value;
}

/**
  * @param {string} message
*/
function appendToGameLogs(message) {
  document.getElementById("events-log").innerHTML += "\n<div>" + message + "</div>";
}

/**
  * @param {string} newRoom
  **/
function changeRoom(newRoom) {
  const currentRoom = document.getElementById("room-name");
  if (newRoom === "" || currentRoom == newRoom) {
    return
  }
  const payload = new ChangeRoomEvent(newRoom);
  sendWsEvent(eventTypes.changeRoom, payload)
}


/**
  * @param {WsEvent} data
  **/
function routeEvent({ type, payload }) {
  console.log("Received Event:", type);
  if (type === eventTypes.receiveImpact) {
    handleReceiveInpact(payload)
    return
  } else {
    appendToGameLogs("Received unhandled message: " + type);
  }
}


/**
  * @param {Object} payload
**/
function handleReceiveInpact(payload) {
  const impactData = Object.assign(new ReceiveImpactEvent, payload)
  const logMessage = `Received impact from: ${impactData}`;
  appendToGameLogs(logMessage);
}

/**
  * @param {string} target
 */
function sendShootTarget(target) {
  if (!target) {
    appendToGameLogs("Shoot: Selected none.");
    return;
  } else if (target == "any") {
    appendToGameLogs("Shoot: <b>Any</b> is not yet available.");
    return;
  }

  const payload = new SendShootEvent(target)
  sendWsEvent(eventTypes.shoot, payload);
}

/**
  * @param {string} data
  */
function handleWsMessage(data) {
  const eventData = JSON.parse(data);
  const event = Object.assign(new WsEvent, eventData)
  routeEvent(event)
}

/**
  * @param {string} eventType
  * @param {Object} payload
 */
function sendWsEvent(eventType, payload) {
  const message = new WsEvent(eventType, payload)
  conn.send(message.toString())
}
