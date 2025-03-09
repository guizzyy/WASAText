<script>
import {RouterLink} from "vue-router";

export default {
  components: RouterLink,
  data: function() {
    return {
      error: null,
      myID: sessionStorage.getItem("ID"),
      myUsername: sessionStorage.getItem("username"),
      myPhoto: sessionStorage.getItem("photo"),
      convID: this.$route.params.convID,
      messages: [],
      sentMessage: "",
      sentPhoto: "",
      receiverID: "",
      receiverName: "",
      receiverPhoto: "",
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
        this.messages = Array.isArray(response.data.reverse()) ? response.data : [];
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

  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6">WASA Text</a>
    <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
  </header>

  <div class="main-container">
    <div>
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
          <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>Options</span>
          </h6>
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/conversations" class="nav-link" @click="getConversation">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                Homepage
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink :to=" '/users/' + this.myID + '/' " class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
                Profile
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>

      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <div class="home-messages">
          <h1 v-if="this.messages.length === 0">
            No messages sent yet...
          </h1>
          <div v-else class="chat-box">
            <div class="messages-list">
              <div v-for="mess in messages" :key="mess.id" class="message">
                {{ mess.text }}
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
  text-align: center;
  margin-top: 20px;
  padding: 20px;
  border-radius: 8px;
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

.message {
  color: white;
  position: relative;
  text-align: center;
}

.messages-list {
  height: auto;
  overflow-y: scroll;
}

</style>