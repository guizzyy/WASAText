<script>

export default {
  data: function () {
    return {
      emojis: [],
      emojiOptions: ["😂", "👍", "🔥", "😡", "😭"],
      showEmojiList: false,
      showEmojiSelect: false,
      showMessage: false,
      showChat: false,
      showMenu: false,
    }
  },

  props: {
    message: Object,
    myID: Number,
  },

  emits: ["updateShowChat", "updateReplyMessage"],

  computed: {
    formattedTimestamp() {
      return new Date(this.message.timestamp).toLocaleTimeString("it-IT", { hour: "2-digit", minute: "2-digit" });
    }
  },

  mounted() {
    this.getComments();
  },

  methods: {
    toggleEmojiSelect()  {
      this.showEmojiSelect = !this.showEmojiSelect;
      if (this.showMenu) this.showMenu = false;
    },
    toggleMenu() {
      this.showMenu = !this.showMenu;

    },
    toggleReactionList() {
      this.showEmojiList = !this.showEmojiList;
      if (this.showMenu) this.showMenu = false;
    },
    toggleMessage() {
      this.showMessage = !this.showMessage;
      if (this.showMenu) this.showMenu = false;
    },
    toggleChatsSelect() {
      this.showChat = !this.showChat;
      this.$emit("updateShowChat", {showChat: this.showChat, messID: this.message.id});
      if (this.showMenu) this.showMenu = false;
    },

    toggleReplySelect(){
      this.$emit("updateReplyMessage", {message: this.message})
      if (this.showMenu) this.showMenu = false;
    },

    async deleteMessage() {
      try {
        await this.$axios.delete(`conversations/${this.message.conv}/messages/${this.message.id}`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.$router.go(0);
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response.data;
        } else if (e.response?.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    },

    async getComments() {
      try {
        let response = await this.$axios.get(`/conversations/${this.message.conv}/messages/${this.message.id}/reactions`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.emojis = response.data;
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response.data;
        } else if (e.response?.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    },

    async commentMessage(emoji) {
      try {
        this.error = null;
        await this.$axios.put(`/conversations/${this.message.conv}/messages/${this.message.id}/reactions`, {emoji: emoji}, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        await this.getComments();
        this.showEmojiSelect = false;
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response.data;
        } else if (e.response?.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    },

    async uncommentMessage() {
      this.error = null;
      try {
        await this.$axios.delete(`conversations/${this.message.conv}/messages/${this.message.id}/reactions`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        await this.getComments();
        this.toggleReactionList();
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response.data;
        } else if (e.response?.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    }
  },
}
</script>

<template>
  <div :class="{'my-mess': message.sender.id === myID, 'receiver-mess': message.sender.id !== myID}" class="mess-wrapper">
    <div class="mess-bubble">
      <div v-if="message.is_forwarded" style="font-size: 10px; color: gray"> forwarded </div>
      <div v-if="message.sender.id !== myID">
        <strong>{{message.sender.username}}</strong>
      </div>
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

      <div class="menu-container position-absolute" style="top: 5px; right: 5px;">
        <i class="menu-icon fas fa-ellipsis-v" @click="toggleMenu" style="cursor: pointer; color: white; font-size: 1.2rem;"></i>
        <div v-if="showMenu" class="menu-popup">
          <i class="action-icon fas fa-mail-reply" @click="toggleReplySelect" title="Reply"></i>
          <i class="action-icon fas fa-mail-forward" @click="toggleChatsSelect" title="Forward"></i>
          <i v-if="message.sender.id === myID" class="action-icon fas fa-trash" @click="toggleMessage" title="Delete"></i>
          <i class="action-icon fas fa-angry" @click="toggleEmojiSelect" title="React"></i>
        </div>
      </div>

      <div v-if="emojis && emojis.length > 0" class="reaction-box" :class="{'receiver-reaction': message.sender.id !== myID}" @click="toggleReactionList">
        <span v-for="(emoji, index) in emojis.slice(0, 3)" :key="index" class="reaction-icon">
          {{ emoji.emoji }}
        </span>
        <span v-if="emojis.length > 3" class="more-reactions">+{{ emojis.length - 3 }}</span>
      </div>

      <div v-if="showEmojiList" class="reaction-popup" :class="{'receiver-popup': message.sender.id !== myID}">
        <div class="reaction-header">
          <strong>Reactions</strong>
          <button class="close-btn" @click="showEmojiList = false">✖</button>
        </div>
        <div class="reaction-content">
          <div v-for="emoji in emojis" :key="emoji.emoji" class="reaction-item">
            {{ emoji.emoji }} - {{ emoji.user.username }}
            <button
                v-if="emoji.user.id === myID"
                class="delete-reaction-btn"
                @click="uncommentMessage"
            >
              ❌
            </button>
          </div>
        </div>
      </div>

      <div v-if="showEmojiSelect" class="emoji-picker">
        <span v-for="emoji in emojiOptions" :key="emoji" class="emoji-choice" @click="commentMessage(emoji)">
          {{ emoji }}
        </span>
      </div>
    </div>

    <div v-if="showMessage" class="overlay">
      <div class="search-box text-black">
        <strong>Are you sure you want to delete the message?</strong>
        <button @click="deleteMessage"> Yes </button>
        <button @click="toggleMessage"> No </button>
      </div>
    </div>

    <div v-if="showChat" class="overlay" @click="toggleChatsSelect">
      <div class="justify-content-center align-items-center">
        <i class="fas fa-arrow-left"></i>
        <strong class="w-50"> Choose where to forward the message </strong>
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

.action-icon {
  cursor: pointer;
  font-size: 1.5rem;
  color: white;
  border-radius: 50%;
  padding: 4px;
  display: inline-block;
}

.mess-photo {
  max-width: 50%;
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

.menu-popup {
  position: absolute;
  top: 25px;
  right: 0;
  background: rgba(0, 0, 0, 0.8);
  border-radius: 5px;
  padding: 5px;
  display: flex;
  flex-direction: column;
  gap: 5px;
  z-index: 10;
}

.menu-popup .action-icon {
  font-size: 1.2rem;
  color: white;
  cursor: pointer;
}

.emoji-picker {
  position: absolute;
  bottom: 110%;
  right: 0;
  background: white;
  border-radius: 10px;
  box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.2);
  padding: 10px;
  display: flex;
  gap: 10px;
  z-index: 10;
}

.emoji-choice {
  font-size: 1.5rem;
  cursor: pointer;
  transition: transform 0.2s;
}

.emoji-choice:hover {
  transform: scale(1.2);
}

.reaction-box {
  position: absolute;
  bottom: -5px;
  left: 10px;
  background: darkgray;
  padding: 2px 8px;
  border-radius: 15px;
  cursor: pointer;
  display: flex;
  gap: 5px;
  align-items: center;
  font-size: 1rem;
}
.receiver-reaction {
  left: auto;
  right: 10px;
}

.reaction-popup {
  position: absolute;
  bottom: 30px;
  left: 10px;
  width: 220px;
  background: white;
  border-radius: 10px;
  padding: 10px;
  z-index: 10;
  color: black;
}
.receiver-popup {
  left: auto;
  right: 10px;
}
.reaction-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  margin-bottom: 5px;
}
.close-btn {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
}
.reaction-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.reaction-item {
  font-size: 1.2rem;
  padding: 5px;
  border-radius: 5px;
}
.reaction-item:hover {
  background: rgba(0, 0, 0, 0.1);
}
.delete-reaction-btn {
  margin-left: 10px;
  border: none;
  background: none;
  cursor: pointer;
  color: red;
  font-size: 1rem;
}
.delete-reaction-btn:hover {
  opacity: 0.7;
}
</style>