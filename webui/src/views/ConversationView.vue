<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  components: {RouterLink, ErrorMsg},
  data: function() {
    return {
      error: null,
      myID: parseInt(sessionStorage.getItem("ID")),
      myUsername: sessionStorage.getItem("username"),
      myPhoto: sessionStorage.getItem("photo") || "https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg",
      convID: this.$route.params.convID,
      messages: [],
      convs: [],
      sentMessage: "",
      sentPhoto: "",
      membersConv: [],

      showLoading: false,

    }
  },

  mounted() {
    this.getConversation()
  },

  methods: {
    scrollToBottom() {
      this.$nextTick( () => {
        const chatBox = document.querySelector(".messages-list");
        if (chatBox) {
          chatBox.scrollTop = chatBox.scrollHeight;
        }
      })
    },

    async sendMessage(mess) {
      if (mess.length === 0) {
        this.error = "Can't send an empty message";
      }
      try {
        let formData = new FormData();
        formData.append('photo', this.sentPhoto);
        formData.append('text', this.sentMessage);
        let response = await this.$axios.post(`conversations/${this.convID}/messages`, formData, {
          headers: {
            Authorization: sessionStorage.getItem("ID"),
            "Content-type": "multipart/form-data"
          }
        });
        this.messages.push(response.data);
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

    async forwardMessage() {

    },

    async deleteMessage() {

    },

    async getConversation() {
      this.showLoading = true;
      this.error = null;
      try {
        let response = await this.$axios.get(`/conversations/${this.convID}/open`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.messages = response.data.messages.reverse();
        this.membersConv = response.data.members;
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
    }
  }
}
</script>

<template>

  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-1 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-5">WASA Text</a>

    <div class="set-buttons d-flex align-items-center me-3">
      <button class="icon-btn" aria-label="Home">
        <router-link to="/conversations" class="icon-btn">
          Home
        </router-link>
      </button>
      <button class="icon-btn" aria-label="Profile">
        <router-link :to="'/users/' + this.myID" class="icon-btn">
          Profile
        </router-link>
      </button>
      <button class="icon-btn" aria-label="Logout">
        <router-link to="/" class="icon-btn">
          Logout
        </router-link>
      </button>
      <div>
        <img :src="this.myPhoto" alt="Stored image" class="profile-pic-header">
      </div>
    </div>
  </header>

  <div class="d-flex position-relative">
    <div class="d-flex position-absolute top-0 end-0 mt-3" style="padding-right: 10px">
      <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
    </div>
  </div>

  <div class="main-container">
    <div>
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
      </nav>

      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <div class="home-messages">
          <h1 v-if="this.messages.length === 0">
            No messages sent yet...
          </h1>
          <div v-else class="chat-box">
            <div class="messages-list">
              <div v-for="mess in messages" :key="mess.id" :class="{'my-mess': mess.sender === myID, 'receiver-mess': mess.sender !== myID}">
                <div class="mess-bubble">
                  <div v-if="mess.text"> {{ mess.text }} </div>
                  <div v-if="mess.photo">
                    <img :src="mess.photo" alt="Message photo" class="mess-photo">
                  </div>
                  <div class="mess-info">
                    <span> {{ mess.timestamp }} </span>
                    <span class="status" v-if="mess.sender === myID">
                      <template v-if="mess.status === 'Read'">
                        <i class="check-mark read">✔✔</i>
                      </template>
                      <template v-if="mess.status === 'Received'">
                        <i class="check-mark read">✔</i>
                      </template>
                    </span>
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
        </div>
      </main>
    </div>
  </div>
</template>

<style>

.home-messages {
  display: flex;
  margin-top: 10px;
  flex-direction: column;
  overflow: hidden;
}

.chat-input-box {
  height: 10%;
  position: absolute;
  bottom: 1em;
  width: 80%;
  padding: 10px;
  display: flex;
  box-shadow: 0 -2px 5px rgba(0, 0, 0, 0.1);
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

/* Optionally, color the check marks differently */
.check-mark.received {
  color: #ccc;
}
.check-mark.read {
  color: #4caf50;
}

</style>