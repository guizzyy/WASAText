<script>
import {RouterLink} from "vue-router";

export default {
  components: RouterLink,
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
        this.messages = response.data.messages;
        this.membersConv = response.data.members;
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
                <div class="mess-bubble"> {{ mess.text }} </div>
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
</template>

<style>

.home-messages {
  display: flex;
  margin-top: 20px;
  border-radius: 8px;
  flex-direction: column;
}

.chat-input-box {
  height: 10%;
  position: absolute;
  bottom: 1em;
  left: 19em;
  width: 80%;
  padding: 10px;
  display: flex;
  align-items: center;
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
  height: fit-content;
}

.messages-list {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  height: 100%;
  overflow-y: auto;
  padding: 20px;
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

</style>