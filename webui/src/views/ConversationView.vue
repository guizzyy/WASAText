<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";
import MessageItem from "../components/MessageItem.vue";
import NotificationMsg from "../components/NotificationMsg.vue";

export default {
  components: {RouterLink, ErrorMsg, MessageItem, NotificationMsg},
  data: function() {
    return {
      error: null,
      myID: parseInt(sessionStorage.getItem("ID")),
      myUsername: sessionStorage.getItem("username"),
      myPhoto: sessionStorage.getItem("photo") || "https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg",
      barConvs: {},
      currConvID: this.$route.params.convID,
      currConv: {},
      allConvMessages: {},
      sentMessage: "",
      sentPhoto: null,
      replyTo: null,

      polling: null,
      report: "",
      showChat: false,
      sentPhotoPreview: null,
      selectedFile: null,

    }
  },

  computed: {
    sortedConvs() {
      return Object.values(this.barConvs).sort((a, b) => {
        return new Date(b.last_message.timestamp) - new Date(a.last_message.timestamp);
      });
    },
  },

  watch : {
    "$route": {
      immediate: true,
      handler(to) {
        this.currConvID = to.params.convID;
        this.getConversation(this.currConvID)
      },
    },
  },

  mounted() {
    this.getConversation(this.currConvID);
    this.startPolling();
  },
  beforeUnmount() {
    this.stopPolling();
  },

  methods: {
    logout() {
      sessionStorage.clear();
      this.$router.push({path: "/"});
    },

    startPolling() {
      this.polling = setInterval(() => {
        this.getConversation(this.currConvID);
      }, 5000);
    },
    stopPolling() {
      clearInterval(this.polling);
    },

    handleReply(mess) {
      this.replyTo = mess;
    },

    handleForwardMessage({feedback}) {
      this.report = feedback;
      setTimeout(() => {
        this.report = "";
      }, 3000);
    },

    isNewDay(index) {
      if (index === 0) return true;
      let prevMess = this.allConvMessages[this.currConvID][index - 1];
      let currMess = this.allConvMessages[this.currConvID][index];
      return new Date(prevMess.timestamp).toDateString() !== new Date(currMess.timestamp).toDateString()
    },

    onFileChange(event) {
      let file = event.target.files[0];
      if (file) {
        this.selectedFile = file;
        this.sentPhoto = file;
        this.sentPhotoPreview = URL.createObjectURL(file)
      }
    },

    removeSelectedFile() {
      this.selectedFile = null;
      this.sentPhoto = null;
      this.sentPhotoPreview = null;
    },

    scrollToBottom() {
      this.$nextTick( () => {
        const chatBox = document.querySelector(".messages-list");
        if (chatBox) {
          chatBox.scrollTop = chatBox.scrollHeight;
        }
      })
    },

    async sendMessage() {
      if (!this.sentMessage && !this.sentPhoto) {
        this.error = "Can't send an empty message";
      }
      try {
        let formData = new FormData();
        if (this.sentPhoto) {
          formData.append('photo', this.sentPhoto);
        }
        if (this.sentMessage) {
          formData.append('text', this.sentMessage);
        }
        if (this.replyTo) {
          formData.append('reply', this.replyTo.message.id);
        }
        let response = await this.$axios.post(`conversations/${this.currConvID}/messages`, formData, {
          headers: {
            Authorization: sessionStorage.getItem("ID"),
            "Content-type": "multipart/form-data"
          }
        });
        this.allConvMessages[this.currConvID].push(response.data);
        this.sentMessage = "";
        this.sentPhoto = "";
        this.replyTo = null;
        this.removeSelectedFile();
        this.scrollToBottom();
        this.barConvs = (await this.$axios.get(`/conversations`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        })).data;
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

    async getConversation(convID) {
      this.error = null;
      try {
        if (!this.allConvMessages[convID]) {
          this.allConvMessages[convID] = [];
        }
        let response = await this.$axios.get(`/conversations/${convID}/open`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.currConv = { ...response.data };
        this.allConvMessages[this.currConvID] = Array.isArray(response.data.messages) ? response.data.messages.reverse() : []
        this.barConvs = (await this.$axios.get(`/conversations`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        })).data;
        this.scrollToBottom();
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
            <ErrorMsg v-if="error" :msg="error" />
            <NotificationMsg v-if="report" :message="report" />
          </div>
        </div>

        <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
          <div v-if="!barConvs || Object.keys(barConvs).length === 0" class="d-flex justify-content-center align-items-center text-center">
            <p class="text-black">No conversation started yet...</p>
          </div>
          <div v-else class="chat-list h-100 d-flex flex-column">
            <router-link v-for="(conv, index) in sortedConvs" :key="index" :to="'/conversations/' + conv.id" class="chat-item d-flex align-items-center p-2">
              <img :src="conv.conv_photo || 'https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg'" alt="Conv photo" class="rounded-circle flex-shrink-0" width="50" height="50">
              <span class="ms-3">{{ conv.name }}</span>
            </router-link>
          </div>
        </nav>

        <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4 position-relative">
          <div v-if="currConv" class="receiver-bar d-flex align-items-center px-3" style="z-index: 30">
            <img :src="currConv.photo || 'https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg'" alt="Conv Photo" class="rounded-circle me-3" width="50" height="50">
            <router-link v-if="currConv.type === 'group'" :to="'/conversations/' + currConvID + '/manage'" class="text-white text-decoration-none ">
              <strong> {{ currConv.name }} </strong>
            </router-link>
            <strong v-else class="text-white">{{ currConv.name }}</strong>
          </div>

          <div class="home-messages">
            <h1 v-if="allConvMessages[currConvID] === 0">No messages sent yet...</h1>

            <div v-else class="chat-box">
              <div class="messages-list">
                <template v-for="(mess, index) in allConvMessages[currConvID]" :key="mess.id">
                  <div v-if="isNewDay(index)" class="text-lg-center fw-bold" style="color: gray; font-size: 14px; margin: 10px 0; display: flex; justify-content: center">
                    <span class="bg-white" style="padding: 5px 10px; border-radius: 10px"> {{ new Date(mess.timestamp).toLocaleDateString("it-IT", {weekday: 'long', month: 'long', day: 'numeric'}) }} </span>
                  </div>
                  <MessageItem :message="mess" :my-i-d="myID" @update-reply-message="handleReply" @update-forward="handleForwardMessage" />
                </template>
              </div>
            </div>
            <div class="chat-input-container">
              <div class="chat-input-box">
                <div v-if="sentPhotoPreview" class="photo-preview d-flex align-items-center" style="gap: 10px; background: #f8f9fa; padding: 8px; border-radius: 10px; max-width: 250px">
                  <img :src="sentPhotoPreview" alt="Photo Preview" class="preview-image">
                  <button class="remove-btn bg-danger text-white border-0" style="padding: 5px; border-radius: 5px; cursor: pointer;" @click="removeSelectedFile">Remove</button>
                </div>
                <div class="w-100 position-relative">
                  <div v-if="replyTo" class="bg-white d-flex justify-content-sm-between align-items-center position-absolute bottom-100" style="padding: 8px; border-radius: 10px 10px 0 0; font-size: 14px; left: 20px">
                    <div class="d-flex align-items-center" style="gap: 5px;">
                      <strong style="color: black"> Reply to: </strong>
                      <span v-if="replyTo.message.photo" class="fw-bold"> 📷 </span>
                      <span v-if="replyTo.message.text" class="text-black">{{ replyTo.message.text }}</span>
                    </div>
                    <button class="remove-reply-btn" style="background: none; border: none; cursor: pointer; font-size: 16px; color: red" @click="replyTo = null">✖</button>
                  </div>
                  <input v-model="sentMessage" type="text" placeholder="Type a message..." class="message-input w-100" maxlength="250" @keyup.enter="sendMessage">
                  <div class="position-absolute d-flex align-items-center cursor-pointer text-secondary attachment">
                    <input type="file" accept="image/*" class="position-absolute w-100 h-100 file-input" @change="onFileChange">
                    <i class="fas fa-paperclip paper-clip" />
                  </div>
                </div>
                <button class="send-button" @click="sendMessage">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="22" y1="2" x2="11" y2="13" />
                    <polygon points="22 2 15 22 11 13 2 9 22 2" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>


<style>

.home-messages {
  display: flex;
  position: relative;
  flex-direction: column;
  overflow: hidden;
  padding-top: 70px;
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
  height: calc(100% - 64px);
  flex-direction: column;
  overflow-y: auto;
  padding: 20px;
  max-height: 98%;
}

.chat-input-container {
  width: 100%;
  position: absolute;
  bottom: -0.5em;
  right: 0.25rem;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.chat-input-box {
  display: flex;
  align-items: center;
  border-radius: 20px;
  padding: 10px;
  position: relative;
}

.message-input {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 20px;
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

.attachment {
  right: 10px;
  top: 0;
  height: 100%;
}

.file-input {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.attachment-wrapper i {
  font-size: 1.2rem;
  padding: 5px;
  cursor: pointer;
  transition: color 0.3s;
}

.attachment-wrapper:hover i {
  color: black;
}

.photo-preview {
  display: flex;
  gap: 10px;
  background: #f8f9fa;
  padding: 8px;
  border-radius: 10px;
}

.preview-image {
  width: 50px;
  height: 50px;
  border-radius: 10px;
}

.remove-btn {
  background: red;
  color: white;
  border: none;
  padding: 5px;
  border-radius: 5px;
  cursor: pointer;
}

.paper-clip{
  font-size: 21px;
}

</style>