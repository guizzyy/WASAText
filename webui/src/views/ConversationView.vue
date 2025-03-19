<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";
import NotificationMsg from "../components/NotificationMsg.vue";

export default {
  components: {NotificationMsg, RouterLink, ErrorMsg},
  data: function() {
    return {
      error: null,
      myID: parseInt(sessionStorage.getItem("ID")),
      myUsername: sessionStorage.getItem("username"),
      myPhoto: sessionStorage.getItem("photo") || "https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg",
      myConvs: JSON.parse(sessionStorage.getItem("convs")) || [],
      currConvID: this.$route.params.convID,
      currentConv: null,
      sentMessage: "",
      sentPhoto: "",
      reactionOf: null,
      emojis: [],
      destinationConv: null,

      showLoading: false,
    }
  },

  mounted() {
    this.getConversation(this.currConvID)
  },

  watch : {
    "$route": {
      immediate: true,
      handler(to) {
        this.currConvID = to.params.convID
        this.getConversation(this.currConvID)
      },
    },
  },

  methods: {
    logout() {
      sessionStorage.clear();
      this.$router.push({path: "/"});
    },

    scrollToBottom() {
      this.$nextTick( () => {
        const chatBox = document.querySelector(".messages-list");
        if (chatBox) {
          chatBox.scrollTop = chatBox.scrollHeight;
        }
      })
    },

    toggleReactions(messID) {
      this.reactionOf = this.reactionOf === messID ? null: messID;
    },

    async commentMessage() {

    },

    async forwardMessage(messID) {
      try {
        this.error = null;
        let response = this.$axios.post(`/conversations/${this.currConvID}/messages/${messID}`,
            {id: this.destinationConv},
            {headers: {Authorization: sessionStorage.getItem("ID")}}
        );
        this.$router.push({path: `/conversations/${this.destinationConv}/open`});
        // to continue 
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

    async deleteMessage(messID) {
      try {
        this.error = null;
        await this.$axios.delete(`/conversations/${this.currConvID}/messages/${messID}`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        })
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

    async sendMessage(mess) {
      if (mess.length === 0) {
        this.error = "Can't send an empty message";
      }
      try {
        let formData = new FormData();
        formData.append('photo', this.sentPhoto);
        formData.append('text', this.sentMessage);
        let response = await this.$axios.post(`conversations/${this.currConvID}/messages`, formData, {
          headers: {
            Authorization: sessionStorage.getItem("ID"),
            "Content-type": "multipart/form-data"
          }
        });
        this.currentConv.messages.push(response.data);
        this.sentMessage = "";
        this.sentPhoto = "";
        this.scrollToBottom();
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response;
        } else if (e.response?.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      } finally {
        this.showLoading = false;
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
     },

    async getConversation(convID) {
      this.showLoading = true;
      this.error = null;
      try {
        let response = await this.$axios.get(`/conversations/${convID}/open`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.currentConv = { ...response.data, messages: Array.isArray(response.data.messages) ? [...response.data.messages].reverse() : [] };
        sessionStorage.setItem("currentConv", JSON.stringify(this.currentConv));
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response.data;
        } else if (e.response?.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      } finally {
        this.showLoading = false;
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    },
  }
}
</script>

<template>
  <div>
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-1 shadow">
      <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-5">WASA Text</a>
      <div class="set-buttons d-flex align-items-center me-3 ms-auto gap-3">
        <button class="icon-btn" aria-label="Home">
          <router-link to="/conversations" class="icon-btn">Home</router-link>
        </button>
        <button class="icon-btn" aria-label="Profile">
          <router-link :to="'/users/' + myID" class="icon-btn">Profile</router-link>
        </button>
        <button class="icon-btn" aria-label="Logout" @click="logout">Logout</button>
        <div>
          <img :src="myPhoto" alt="Stored image" class="profile-pic-header">
        </div>
      </div>
    </header>

    <div class="container-fluid">
      <div class="row">
        <div class="d-flex position-relative">
          <div class="d-flex position-absolute top-0 end-0 mt-3">
            <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
          </div>
        </div>

        <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
          <div v-if="!myConvs || myConvs.length === 0" class="h-100 mt-3 d-flex justify-content-center align-items-center text-center">
            <p class="text-black">No conversation started yet...</p>
          </div>
          <div v-else class="chat-list h-100 mt-2 d-flex flex-column">
            <router-link v-for="conv in myConvs" :key="conv.id" :to="'/conversations/' + conv.id" class="chat-item d-flex align-items-center p-2">
              <img :src="conv.photo || 'https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg'" alt="Conv photo" class="rounded-circle flex-shrink-0" width="50" height="50">
              <span class="ms-3">{{ conv.name }}</span>
            </router-link>
          </div>
        </nav>

        <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4 position-relative">
          <div v-if="currentConv" class="receiver-bar d-flex align-items-center px-3">
            <img :src="currentConv.photo || 'https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg'" alt="Conv Photo" class="rounded-circle me-3" width="50" height="50">
            <router-link v-if="currentConv.type === 'group'" :to="'/conversations/' + currConvID + '/manage'" class="text-white text-decoration-none ">
              <strong> {{ currentConv.name }} </strong>
            </router-link>
            <strong v-else class="text-white">{{ currentConv.name }}</strong>
          </div>

          <div class="home-messages">
            <h1 v-if="!currentConv || currentConv.messages.length === 0">No messages sent yet...</h1>

            <div v-else class="chat-box">
              <div class="messages-list">

                <div v-for="mess in currentConv.messages" :key="mess.id" :class="{'my-mess': mess.sender.id === myID, 'receiver-mess': mess.sender.id !== myID}" class="mess-wrapper">
                  <div class="mess-bubble">
                    <div v-if="mess.text">{{ mess.text }}</div>
                    <div v-if="mess.photo">
                      <img :src="mess.photo" alt="Message photo" class="mess-photo">
                    </div>

                    <div class="mess-info">
                      <span>{{ new Date(mess.timestamp).toLocaleDateString("it-IT", {hour: "numeric", minute: "numeric"}) }}</span>

                      <span class="status" v-if="mess.sender === myID">
                        <template v-if="mess.status === 'Read'">
                          <i class="check-mark read">✔✔</i>
                        </template>
                        <template v-else-if="mess.status === 'Received'">
                          <i class="check-mark received">✔</i>
                        </template>
                      </span>
                    </div>

                    <div class="mess-actions" :class="{'my-actions': mess.sender.id === myID, 'receiver-actions': mess.sender.id !== myID}">
                      <i class="action-icon bi bi-forward" @click="forwardMessage"></i>
                      <i v-if="mess.sender.id === myID" class="action-icon fa-solid fa-delete-left" @click="deleteMessage"></i>
                      <i class="action-icon fa-brands fa-react" @click="toggleReactions"></i>
                      <div v-if="reactionOf === mess.id" class="emoji-list">
                        <span v-for="emoji in emojis" :key="emoji" @click="commentMessage">{{ emoji }}</span>
                      </div>
                    </div>

                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="chat-input-box">
            <input v-model="sentMessage" type="text" placeholder="Type a message..." class="message-input" @keyup.enter="sendMessage" maxlength="250">
            <button @click="sendMessage" class="send-button">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="22" y1="2" x2="11" y2="13"></line>
                <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
              </svg>
            </button>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>


<style>

.home-messages {
  display: flex;
  margin-top: 10px;
  flex-direction: column;
  overflow: hidden;
  padding-top: 50px;
}

.chat-box{
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 160px);
  overflow: hidden;
}

.messages-list {
  display: flex;
  justify-content: flex-start;
  height: auto;
  flex-direction: column;
  overflow-y: auto;
  padding: 20px;
  max-height: 98%;
}

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
  font-size: 1.2rem;
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

.chat-input-box {
  height: 10%;
  justify-content: flex-end;
  position: fixed;
  bottom: 1em;
  right: 0;
  width: 83%;
  padding: 10px;
  display: flex;
}

.message-input {
  height: 100%;
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 20px;
  outline: none;
}

.send-button {
  height: 100%;
  margin-left: 10px;
  padding: 10px 15px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  transition: background 0.3s;
}

.send-button:hover {
  background-color: #0056b3;
}

.chat-item {
  align-items: center;
  padding: 2vh;
  display: block;
  cursor: pointer;
  text-decoration: none;
  color: black;
  border-bottom: 1px solid #ddd;
  background-color: white;
}

.chat-item:hover {
  background-color: lightgray;
}

.receiver-bar {
  width: 100%;
  background-color: #343a40;
  position: absolute;
  left: 0;
  display: flex;
  align-items: center;
  padding: 10px;
}

</style>