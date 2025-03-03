<script>
import {RouterLink} from "vue-router";

export default {
  components: RouterLink,
  data: function() {
    return {
      error: null,
      ID: sessionStorage.getItem("ID"),
      username: sessionStorage.getItem("username"),
      photo: sessionStorage.getItem("photo"),
      message: sessionStorage.getItem("message"),
      convs: [],
      newUser: "",
      searchResults: [],
      userSelected: null,
      newConv: {},

      showLoading: false,
      showUserSearch: false,
    }
  },

  mounted() {
    this.getConversations();
  },

  methods: {
    openSearchBar() {
      this.showUserSearch = true;
    },

    closeSearchBar() {
      this.showUserSearch = false;
      this.newUser = "";
      this.searchResults = [];
      this.userSelected = null;
    },

    async searchUsers() {
      this.error = null;
      this.showLoading = true;
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
          this.error = e.response;
        } else if (e.response && e.response.status === 500) {
          this.error = e.response.data
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    },

    async getConversations() {
      this.error = null;
      this.showLoading = true;
      try {
        let response = await this.$axios.get("/conversations", {
          headers: {
            Authorization: "Bearer " + sessionStorage.getItem("ID"),
          }
        });
        this.convs = response.data;
      } catch (e) {
        this.showLoading = false;
        if (e.response?.status === 400) {
          this.error = "Failed to get conversations.";
        } else if (e.response?.status === 404) {
          this.error = "User Not Found";
        } else if (e.response?.status === 500) {
          this.error = "Server Error, please try again";
          console.error("Error fetching the conversations")
        } else {
          this.error = e.toString();
        }
        setTimeout(() => {
          this.error = null;
        }, 3000)
      } finally {
        this.showLoading = false;
      }
    },

    async getConversation(convID) {
      this.error = null;
      this.showLoading = true;
      let response = await this.$axios.get("/conversations/:convID", {
        headers: {
          Authorization: sessionStorage.getItem("ID")
        }
      })
    },

    async startConversation(user) {
      this.error = null;
      this.showLoading = true;
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
        this.$router.push({path: `/conversations/${this.newConv.id}`});
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = "Invalid username (it must be between 3 and 16 characters).";
        } else if (e.response?.status === 500) {
          this.error = "Server Error, please try again later.";
        } else {
          this.error = "An unexpected error occurred.";
          console.error(e); // Log for debugging
        }
      } finally {
        this.closeSearchBar();
        this.showLoading = false;
      }
      setTimeout(() => {
        this.error = null;
      }, 2500)
    }
  },
}
</script>

<template>

  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA Text</a>
    <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
  </header>

  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
          <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>Options</span>
          </h6>
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/conversations" class="nav-link" @click="getConversations">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                Homepage
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink :to=" '/users/' + ID + '/' " class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
                Profile
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>

      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <div class="d-flex position-relative">
          <div class="d-flex position-absolute top-0 end-0 mt-3">
            <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
          </div>

        </div>

        <div class="home-container">
          <h1> Chats </h1>

          <p v-if="this.convs.length === 0">No conversation started yet...</p>
          <ul v-else>
            <li v-for="conv in convs" :key="conv.id">
              <router-link to="/conversations/:convID">
                {{ conv.name }}
              </router-link>
            </li>
          </ul>

          <div class="new-chat-button" @click="openSearchBar">
            <svg class="feather" width="24" height="24"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
            <span style="font-size: 30px; position: absolute; justify-content: center; font-weight: bold; bottom: .25rem">+</span>
          </div>

          <div v-if="showUserSearch" class="overlay">
            <div class="search-box position-relative">
              <input v-model="newUser" @input="searchUsers" placeholder="Search for a user..." />
              <ul>
                <li v-for="user in searchResults" :key="user.id" @click="startConversation(user)">
                  {{ user.username }}
                </li>
              </ul>
              <button @click="closeSearchBar">Cancel</button>
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
  margin-top: 20px;
  padding: 20px;
  border-radius: 8px;
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

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(5px);
  display: flex;
  justify-content: center;
  align-items: center;
}

.search-box {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 30%;
  text-align: center;
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

</style>