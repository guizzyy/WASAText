<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";
import NotificationMsg from "../components/NotificationMsg.vue";
import MessageItem from "../components/MessageItem.vue";

export default {
  components: {NotificationMsg, RouterLink, ErrorMsg, MessageItem},
  data: function() {
    return {
      error: null,
      myID: parseInt(sessionStorage.getItem("ID")),
      myUsername: sessionStorage.getItem("username"),
      myPhoto: sessionStorage.getItem("photo") || "https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg",
      myConvs: JSON.parse(sessionStorage.getItem("convs")) || [],
      currConvID: this.$route.params.convID,
      lastMessageIDs: {},
      allConvMessages: {"": []},
      currentConv: {},
      sentMessage: "",
      sentPhoto: null,
      reactionOf: null,
      emojis: [],
      destinationConv: null,

      sentPhotoPreview: null,
      selectedFile: null,

    }
  },

  watch : {
    "$route": {
      immediate: true,
      handler(to) {
        this.currConvID = to.params.convID;
        if (!this.lastMessageIDs[this.currConvID]) {
          this.lastMessageIDs[this.currConvID] = 0;
        }
        this.getConversation(this.currConvID)
      },
    },
  },

  methods: {
    logout() {
      sessionStorage.clear();
      this.$router.push({path: "/"});
    },

    onFileChange(event) {
      let file = event.target.files[0];
      if (file) {
        this.selectedFile = file;
        this.sentPhoto = URL.createObjectURL(file)
      }
    },

    removeSelectedFile() {
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

    async sendMessage() {
      if (!this.sentMessage && !this.sentPhoto) {
        this.error = "Can't send an empty message";
        return
      }
      try {
        let formData = new FormData();
        if (this.sentPhoto) {
          formData.append('photo', this.sentPhoto);
        }
        if (this.sentMessage) {
          formData.append('text', this.sentMessage);
        }
        let response = await this.$axios.post(`conversations/${this.currConvID}/messages`, formData, {
          headers: {
            Authorization: sessionStorage.getItem("ID"),
            "Content-type": "multipart/form-data"
          }
        });
        this.allConvMessages[this.currConvID].push(response.data);
        this.sentMessage = "";
        this.removeSelectedFile();
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
      const lastID = this.lastMessageIDs[convID] || 0;
      try {
        if (!this.allConvMessages[convID]) {
          this.allConvMessages[convID] = [];
        }
        let response = await this.$axios.get(`/conversations/${convID}/open?lastID=${lastID}`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.currentConv = { ...response.data };
        let newMessages = Array.isArray(response.data.messages) ? response.data.messages.reverse() : [];
        if (newMessages.length > 0) {
          this.lastMessageIDs[convID] = newMessages[newMessages.length - 1].id;
          this.allConvMessages[convID].push(...newMessages);
        }
        this.currentConv.messages = this.allConvMessages[convID];
        if (this.currentConv.type === 'group') {
          sessionStorage.setItem("currGroup", JSON.stringify(response.data));
        }
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
              <img :src="conv.conv_photo || 'https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg'" alt="Conv photo" class="rounded-circle flex-shrink-0" width="50" height="50">
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
            <h1 v-if="allConvMessages[currConvID] === 0">No messages sent yet...</h1>

            <div v-else class="chat-box">
              <div class="messages-list">
                <MessageItem
                    v-for="mess in allConvMessages[currConvID]"
                    :key="mess.id"
                    :message="mess"
                    :myID="myID"
                    :reactionOf="reactionOf"
                    :emojis="emojis"
                    @toggle-reactions="toggleReactions"
                    @comment="commentMessage"
                />
              </div>
            </div>
          </div>

          <div class="chat-input-box">
            <input v-model="sentMessage" type="text" placeholder="Type a message..." class="message-input" @keyup.enter="sendMessage" maxlength="250">
            <div class="position-absolute d-flex align-items-center cursor-pointer text-secondary attachment">
              <input type="file" accept="image/*" @click="onFileChange" class="position-absolute w-100 h-100 file-input">
              <i class="fas fa-paperclip fs-1 p-5"></i>
            </div>
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

.chat-input-box {
  height: 10%;
  justify-content: flex-end;
  position: fixed;
  bottom: 1em;
  right: 0;
  width: 83%;
  padding: 10px;
  display: flex;
  align-items: center;
  border-radius: 20px;
}

.message-input {
  height: 100%;
  flex-grow: 1;
  padding: 10px 40px 10px 10px;
  border: 1px solid #ccc;
  border-radius: 20px;
  outline: none;
  position: relative;
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

.attachment{
  position: absolute;
  right: 50px;
  display: flex;
  align-items: center;
  cursor: pointer;
  color: gray;
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

</style>