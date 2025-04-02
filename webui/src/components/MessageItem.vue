<script>

export default {

  props: {
    message: Object,
    myID: Number,
  },

  emits: ["updateReplyMessage", "updateForward"],
  data: function () {
    return {
      emojis: [],
      emojiOptions: ["😂", "👍", "🔥", "😡", "😭"],
      newUser: null,
      searchResults: [],

      showUserSearch: false,
      showEmojiList: false,
      showEmojiSelect: false,
      showMessage: false,
      showMenu: false,
      isForwarded: false,
    }
  },

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
    toggleForward() {
      this.showUserSearch = true;
      if (this.showMenu) this.showMenu = false;
    },
    toggleReplySelect() {
      this.$emit("updateReplyMessage", {message: this.message})
      if (this.showMenu) this.showMenu = false;
    },
    closeSearchBar() {
      this.showUserSearch = false;
    },

    async forwardMessage(username) {
      try {
        this.error = null;
        let response = await this.$axios.post(`/conversations/${this.message.conv}/messages/${this.message.id}`,
            {username: username},
            {headers: {Authorization: sessionStorage.getItem("ID")}}
        );
        let newMess = response.data.message
        let report = response.data.report
        this.showUserSearch = false;
        this.newUser = null;
        this.$emit("updateForward", {feedback: report})
        this.$router.push({path: `/conversations/${newMess.conv}`})
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response;
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
    },

    async searchUsers() {
      clearTimeout(this.searchTimeout);
      this.searchTimeout = setTimeout(async () => {
        this.error = null;
        if (this.newUser.length === 0) {
          this.searchResults = []
        }
        try {
          let response = await this.$axios.get(`/users/${this.myID}/search?username=${this.newUser}`, {
            headers: {
              Authorization: sessionStorage.getItem("ID")
            }
          });
          this.searchResults = response.data;
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
      }, 300);
    },
  },
}
</script>

<template>
  <div :class="{'my-mess': message.sender.id === myID, 'receiver-mess': message.sender.id !== myID}" class="mess-wrapper">
    <div class="mess-bubble">
      <div v-if="message.reply_photo || message.reply_text" class="flex-lg-row">
        <strong> Reply to: </strong>
        <i v-if="message.reply_photo" class="fas fa-camera" />
        <p v-if="message.reply_text"> {{ message.reply_text }} </p>
      </div>
      <div v-if="message.is_forwarded" style="font-size: 10px; color: black"> forwarded </div>
      <div v-if="message.sender.id !== myID">
        <strong>{{ message.sender.username }}</strong>
      </div>
      <div v-if="message.photo">
        <img :src="message.photo" alt="Message photo" class="mess-photo">
      </div>
      <div v-if="message.text">{{ message.text }}</div>

      <div class="mess-info">
        <span>{{ formattedTimestamp }}</span>

        <span v-if="message.sender.id === myID" class="status">
          <template v-if="message.status === 'Read'">
            <i class="check-mark read">✔✔</i>
          </template>
          <template v-else-if="message.status === 'Received'">
            <i class="check-mark received">✔</i>
          </template>
        </span>
      </div>

      <div class="menu-container position-absolute" style="top: 5px; right: 5px;">
        <i class="menu-icon fas fa-ellipsis-v" style="cursor: pointer; color: white; font-size: 1.2rem;" @click="toggleMenu" />
        <div v-if="showMenu" class="menu-popup">
          <i class="action-icon fas fa-mail-reply" title="Reply" @click="toggleReplySelect" />
          <i class="action-icon fas fa-paper-plane" title="Forward" @click="toggleForward" />
          <i v-if="message.sender.id === myID" class="action-icon fas fa-trash" title="Delete" @click="toggleMessage" />
          <i class="action-icon fas fa-angry" title="React" @click="toggleEmojiSelect" />
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

      <div v-if="showEmojiSelect" :class="{'emoji-picker': message.sender.id === myID, 'receiver-emoji-picker': message.sender.id !== myID}">
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

    <div v-if="showUserSearch" class="overlay">
      <div class="search-box position-relative">
        <input v-model="newUser" placeholder="Forward the message to..." @input="searchUsers">
        <ul>
          <li v-for="user in searchResults" :key="user.id" @click="forwardMessage(user.username)">
            {{ user.username }}
          </li>
        </ul>
        <button @click="closeSearchBar">Cancel</button>
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
  padding: 12px 15px;
  border-radius: 15px;
  font-size: 16px;
  position: relative;
  word-wrap: break-word;
  margin-bottom: 15px;
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
  top: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.8);
  border-radius: 5px;
  padding: 5px;
  display: flex;
  flex-direction: row;
  gap: 5px;
  z-index: 10;
  transform: translate(0, -100%);
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
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  padding: 10px;
  display: flex;
  gap: 10px;
  z-index: 10;
}

.receiver-emoji-picker {
  position: absolute;
  bottom: 110%;
  left: 0;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
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
  bottom: 0;
  transform: translate(0, 50%);
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
  left: 10px;
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