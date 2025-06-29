<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";
import NotificationMsg from "../components/NotificationMsg.vue";

export default {
  components: {NotificationMsg, RouterLink, ErrorMsg},
  data: function() {
    return {
      error: null,
      ID: parseInt(sessionStorage.getItem("ID")),
      username: sessionStorage.getItem("username"),
      photo: sessionStorage.getItem("photo") || "https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg",
      report: sessionStorage.getItem("report"),
      convs: {},
      newUser: "",
      searchResults: [],
      newConv: {},
      groupName: "",

      polling: null,
      searchTimeout: null,
      showUserSearch: false,
      showGroupName: false,
      showPopUp: false,
    }
  },

  computed: {
    sortedConvs() {
      return Object.values(this.convs).sort((a, b) => {
        return new Date(b.last_message.timestamp) - new Date(a.last_message.timestamp);
      });
    }
  },

  mounted() {
    this.getConversations();
    setTimeout(() => {
      sessionStorage.removeItem("report");
      this.report = "";
    }, 3000);
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
        this.getConversations();
      }, 5000);
    },
    stopPolling() {
      clearInterval(this.polling);
    },

    togglePopUp() {
      this.showPopUp = !this.showPopUp;
    },

    openGroupNameBar() {
      this.showGroupName = true;
      this.showPopUp = false;
    },

    closeGroupNameBar() {
      this.groupName = "";
      this.showGroupName = false;
    },

    openSearchBar() {
      this.showUserSearch = true;
      this.showPopUp = false;
    },

    closeSearchBar() {
      this.showUserSearch = false;
      this.newUser = "";
      this.searchResults = [];
    },

    async searchUsers() {
      clearTimeout(this.searchTimeout);
      this.searchTimeout = setTimeout(async () => {
        this.error = null;
        if (this.newUser.length === 0) {
          this.searchResults = []
        }
        try {
          let response = await this.$axios.get(`/users/${this.ID}/search?username=${this.newUser}`, {
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

    async getConversations() {
      this.error = null;
      try {
        let response = await this.$axios.get("/conversations", {
          headers: {
            Authorization: sessionStorage.getItem("ID"),
          }
        });
        let conversations = response.data;
        if (!conversations || conversations.length === 0) {
          this.convs = {};
        } else {
          conversations.forEach(conv => {
            this.convs[conv.id] = conv;
          })
        }
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = "Failed to get conversations.";
        } else if (e.response?.status === 404) {
          this.error = "User Not Found";
        } else if (e.response?.status === 500) {
          this.error = "Server Error, please try again";
        } else {
          this.error = e.toString();
        }
        setTimeout(() => {
          this.error = null;
        }, 3000)
      }
    },

    async startConversation(user) {
      this.error = null;
      try {
        let response = await this.$axios.post("/conversations",
            {username: user.username},
            {
              headers: {
                'Content-type': 'application/json',
                Authorization: sessionStorage.getItem("ID"),
              }
            }
        )
        this.newConv = response.data;
        this.convs[this.newConv.id] = this.newConv;
        this.$router.push({path: `/conversations/${this.newConv.id}`});
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = "Invalid username (it must be between 3 and 16 characters).";
        } else if (e.response?.status === 404) {
          this.error = "User not found"
        } else if (e.response?.status === 500) {
          this.error = "Server Error, please try again later.";
        } else {
          this.error = e.response.data
        }
      } finally {
        this.closeSearchBar();
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    },

    async createGroup() {
      this.error = null;
      if (this.groupName.length === 0) {
        this.error = "Write a proper name for a new group";
      }
      else try {
        let response = await this.$axios.post(`/group`, {name: this.groupName}, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.newConv = response.data;
        this.convs[this.newConv.id] = this.newConv;
        this.$router.push({path: `/conversations/${this.newConv.id}`})
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response.data;
        } else if (e.response?.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      } finally {
        this.closeGroupNameBar();
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    }
  },
}
</script>

<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-1 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-5">WASA Text</a>

    <div class="set-buttons d-flex align-items-center me-3 ms-auto gap-3">
      <button class="icon-btn" aria-label="Home">
        <router-link to="/conversations" class="icon-btn">
          Home
        </router-link>
      </button>
      <button class="icon-btn" aria-label="Profile">
        <router-link :to="'/users/' + ID" class="icon-btn">
          Profile
        </router-link>
      </button>
      <button class="icon-btn" aria-label="Logout" @click="logout">
        Logout
      </button>
      <div>
        <img :src="photo" alt="Stored image" class="profile-pic-header">
      </div>
    </div>
  </header>

  <div class="container-fluid">
    <div class="d-flex position-relative">
      <div class="d-flex position-absolute top-0 end-0 mt-3">
        <ErrorMsg v-if="error" :msg="error" />
        <NotificationMsg v-if="report" :message="report" />
      </div>
    </div>

    <div class="row">
      <main style="height: 90%">
        <div class="home-container">
          <h1> Chats </h1>

          <p v-if="!convs || Object.keys(convs).length === 0">No conversation started yet...</p>
          <div v-else class="chat-container">
            <div class="chat-list">
              <router-link v-for="(conv, index) in sortedConvs" :key="index" :to="'/conversations/' + conv.id" class="chat d-flex align-items-start gap-3 p-2">
                <div class="chat-photo h-auto">
                  <img :src="conv.conv_photo || 'https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg' " alt="Conv photo" class="rounded-circle flex-shrink-0" width="50" height="50">
                </div>
                <div class="flex-grow-1">
                  <div class="d-flex justify-content-between">
                    <strong> {{ conv.name }} </strong>
                    <small v-if="conv.last_message.id" class="text-muted"> {{ new Date(conv.last_message.timestamp).toLocaleDateString("it-IT", {hour: "numeric", minute: "numeric"}) }} </small>
                  </div>
                  <div class="d-flex justify-content-between align-items-center">
                    <span class="d-flex align-items-center flex-grow-1 text-truncate">
                      <i v-if="conv.last_message.id && conv.last_message.photo" class="bi bi-camera-fill me-1" />
                      <p class="text-muted text-truncate mb-0">{{ conv.last_message.id ? conv.last_message.text : "No messages sent yet..." }}</p>
                    </span>
                  </div>
                </div>
              </router-link>
            </div>
          </div>

          <div class="new-chat-button" @click="togglePopUp">
            <svg class="feather" width="24" height="24"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
            <span style="font-size: 30px; position: absolute; justify-content: center; font-weight: bold; bottom: .25rem">+</span>
          </div>

          <div v-if="showPopUp" class="chat-popup">
            <button @click="openSearchBar"> Private </button>
            <button @click="openGroupNameBar"> Group </button>
          </div>

          <div v-if="showUserSearch" class="overlay">
            <div class="search-box position-relative">
              <input v-model="newUser" placeholder="Search for a user..." @input="searchUsers">
              <ul>
                <li v-for="user in searchResults" :key="user.id" @click="startConversation(user)">
                  {{ user.username }}
                </li>
              </ul>
              <button @click="closeSearchBar">Cancel</button>
            </div>
          </div>

          <div v-if="showGroupName" class="overlay">
            <div class="search-box position-relative">
              <input v-model="groupName" placeholder="Insert a name for the group...">
              <button class="" @click="createGroup"> Create group </button>
              <button @click="closeGroupNameBar"> Cancel </button>
            </div>
          </div>
        </div>

        <RouterView />
      </main>
    </div>
  </div>
</template>

<style>

.home-container {
  text-align: center;
  padding: 20px;
  border-radius: 8px;
  height: 90%;
}

.new-chat-button {
  position: fixed;
  bottom: 30px;
  right: 30px;
  width: 50px;
  height: 50px;
  background-color: #298dff; /* Change color as needed */
  color: white;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.new-chat-button:hover {
  background-color: #0a53a8;
}

.chat-popup {
  width: 10%;
  position: fixed;
  bottom: 90px;
  right: 30px;
  background: none;
  padding: 10px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.chat-popup button {
  background: #298dff;
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 5px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.chat-popup button:hover {
  background: #0a53a8;
}

.search-box input {
  width: 100%;
  padding: 10px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.search-box ul {
  list-style: none;
  padding: 0;
  max-height: 200px;
  overflow-y: auto;
}

.search-box li {
  padding: 10px;
  cursor: pointer;
  color: white;
  background-color: rgba(0, 0, 255, 0.38);
}

.search-box li:hover {
  background: #0a53a8;
}

.profile-pic-header {
  width: 35px;
  height: 35px;
  border-radius: 50%;
  object-fit: cover;
  background-color: black;
}

.chat-container {
  align-items: center;
  justify-content: center;
  height: auto;
  margin: auto;
  flex-direction: column;
}

.chat-list {
  width: auto;
  height: 80vh;
  overflow-y: scroll;
  border-radius: 10px;
  background-color: #f9f9f9;
}

.chat {
  height: 5em;
  display: block;
  width: 100%;
  text-align: left;
  border: none;
  background-color: white;
  cursor: pointer;
  border-bottom: 1px solid #ddd;
  text-decoration: none;
  color: black;
  font-family: inherit;
}

.chat:hover {
  background-color: lightgray;
}

.set-buttons {
  position: absolute;
  top: auto;
  right: 1em;
  display: flex;
}

.icon-btn {
  display: flex;
  background: none;
  border: none;
  color: white;
  font-size: 1.5rem; /* Adjust size of icons */
  cursor: pointer;
  transition: color 0.3s ease-in-out;
  text-decoration: none;
}

.icon-btn:hover {
  color: #ccc; /* Slight color change on hover */
}

</style>