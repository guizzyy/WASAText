<script>

export default {
  props: {
    message: Object,
    myID: Number,
    reactionOf: Number,
    emojis: Array
  },
  computed: {
    formattedTimestamp() {
      return new Date(this.message.timestamp).toLocaleDateString("it-IT", { hour: "numeric", minute: "numeric" });
    }
  }
}
</script>

<template>
  <div :class="{'my-mess': message.sender.id === myID, 'receiver-mess': message.sender.id !== myID}" class="mess-wrapper">
    <div class="mess-bubble">
      <div v-if="message.photo">
        <img :src="message.photo" alt="Message photo" class="mess-photo">
      </div>
      <div v-if="message.text">{{ message.text }}</div>

      <div class="mess-info">
        <span>{{ formattedTimestamp }}</span>

        <span class="status" v-if="message.sender.id === myID">
          <template v-if="message.status === 'Read'">
            <i class="check-mark read">✔✔</i>
          </template>
          <template v-else-if="message.status === 'Received'">
            <i class="check-mark received">✔</i>
          </template>
        </span>
      </div>

      <div class="mess-actions" :class="{'my-actions': message.sender.id === myID, 'receiver-actions': message.sender.id !== myID}">
        <i class="action-icon bi bi-forward" @click="$emit('forward', message)"></i>
        <i v-if="message.sender.id === myID" class="action-icon fa-solid fa-delete-left" @click="$emit('delete', message.id)"></i>
        <i class="action-icon fa-brands fa-react" @click="$emit('toggle-reactions', message.id)"></i>
        <div v-if="reactionOf === message.id" class="emoji-list">
          <span v-for="emoji in emojis" :key="emoji" @click="$emit('comment', { message, emoji })">{{ emoji }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

.mess-wrapper {
  position: relative;
  margin-bottom: 10px;
}

.my-mess {
  align-self: flex-end;
  max-width: 60%;
}

.receiver-mess {
  align-self: flex-start;
  max-width: 60%;
}

.mess-bubble {
  padding: 10px 15px;
  border-radius: 15px;
  font-size: 16px;
  position: relative;
  word-wrap: break-word;
  margin-bottom: 10px;
}

.my-mess .mess-bubble {
  background-color: #0078ff;
  color: white;
  border-bottom-right-radius: 5px;
}

.my-mess .mess-bubble::after {
  content: "";
  position: absolute;
  right: -10px;
  top: 50%;
  width: 0;
  height: 0;
  border-left: 10px solid #0078ff;
  border-top: 10px solid transparent;
  border-bottom: 10px solid transparent;
  transform: translateY(-50%);
}

.receiver-mess .mess-bubble {
  background-color: #ff9229;
  color: white;
  border-bottom-right-radius: 5px;
}

.receiver-mess .mess-bubble::after {
  content: "";
  position: absolute;
  left: -10px;
  top: 50%;
  width: 0;
  height: 0;
  border-right: 10px solid #ff9229;
  border-top: 10px solid transparent;
  border-bottom: 10px solid transparent;
  transform: translateY(-50%);
}

.mess-actions {
  background-color: white;
  flex-direction: column;
  position: absolute;
  bottom: 5px;
  display: flex;
  gap: 5px;
  opacity: 0;
  transition: opacity 0.2s;
}

.message-wrapper:hover .mess-actions {
  opacity: 1;
}

.my-mess .mess-actions.my-actions {
  left: -50px;
}

.receiver-mess .mess-actions.receiver-actions {
  right: -50px;
}

.action-icon {
  cursor: pointer;
  font-size: 1.5rem;
  background: black;
  color: black;
  border-radius: 50%;
  padding: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
  display: inline-block;
}

.emoji-list {
  position: absolute;
  bottom: 100%;
  right: 0;
  display: flex;
  gap: 5px;
  background: white;
  padding: 5px;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.mess-photo {
  max-width: 100%;
  border-radius: 10px;
  margin-top: 5px;
}

.mess-info {
  font-size: 0.75rem;
  margin-top: 4px;
  display: flex;
  justify-content: flex-end;
  opacity: 0.8;
}

.check-mark {
  font-size: 0.5rem;
  margin-left: 5px;
}

.check-mark.received {
  color: #ccc;
}
.check-mark.read {
  color: #4caf50;
}

</style>