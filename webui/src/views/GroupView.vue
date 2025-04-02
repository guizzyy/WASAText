<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";
import NotificationMsg from "../components/NotificationMsg.vue";

export default {
  components: {NotificationMsg, ErrorMsg, RouterLink},
  data: function() {
    return {
      error: null,
      myID: parseInt(sessionStorage.getItem("ID")),
      myUsername: sessionStorage.getItem("username"),
      myPhoto: sessionStorage.getItem("photo") || "https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg",
      currConvID: parseInt(this.$route.params.convID),
      currGroup: {},
      members: [],
      searchResults: [],
      newMember: "",
      newName: "",
      newPhoto: "",
      selectedFile: null,

      report: "",
      showUserSearch: false,
      showNameBar: false,
      showPhotoBar: false,
    }
  },

  mounted() {
    this.fetchGroup();
  },

  methods: {
    logout() {
      sessionStorage.clear();
      this.$router.push({path: '/'});
    },

    openSearchBar() {
      this.showUserSearch = true;
    },

    closeSearchBar() {
      this.showUserSearch = false;
      this.newMember = "";
      this.searchResults = [];
    },

    openNameBar() {
      this.showNameBar = true;
    },

    openPhotoBar() {
      this.showPhotoBar = true;
    },

    closeNameBar() {
      this.showNameBar = false;
      this.newName = "";
    },

    closePhotoBar() {
      this.showPhotoBar = false;
      this.newPhoto = "";
      this.selectedFile = null;
    },

    onFileChange(event) {
      let file = event.target.files[0];
      if (file) {
        this.selectedFile = file;
        this.newPhoto = URL.createObjectURL(file)
      }
    },

    async setGroupName() {
      try {
        this.error = null;
        let response = await this.$axios.put(`conversations/${this.currConvID}/manage/name`, {name: this.newName},{
          headers: { Authorization: sessionStorage.getItem("ID") }
        });
        this.currGroup.name = this.newName;
        this.report = response.data.report;
        this.closeNameBar();
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response;
        } else if (e.response?.status === 500) {
          this.error = e.response.data;
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
        this.report = "";
      }, 2500)
    },

    async setGroupPhoto() {
      if (!this.selectedFile) {
        this.error = "Please upload a photo";
        return;
      }
      try {
        let formData = new FormData();
        formData.append('photo', this.selectedFile);
        let response = await this.$axios.put(`conversations/${this.currConvID}/manage/photo`, formData, {
          headers : {
            "Content-type" : "multipart/form-data",
            Authorization : sessionStorage.getItem("ID")
          }
        });
        this.report = response.data.report;
        this.currGroup.photo = response.data.photo
        this.closePhotoBar();
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response;
        } else if (e.response?.status === 500) {
          this.error = e.response.data;
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
        this.report = "";
      }, 2500)
    },

    async fetchGroup() {
      try {
        this.error = null;
        let response = await this.$axios.get(`group/${this.currConvID}`, {
          headers: {
            Authorization: sessionStorage.getItem("ID")
          }
        });
        this.currGroup = response.data;
        this.members = response.data.members;
      } catch (e) {
        if (e.response?.status === 400) {
          this.error = e.response;
        } else if (e.response?.status === 500) {
          this.error = e.response.data;
        } else {
          this.error = e.toString();
        }
      }
      setTimeout(() => {
        this.error = null;
        this.report = "";
      }, 2500)
    },

    async addToGroup(user) {
      this.error = null;
      try {
        let response = await this.$axios.post(`/memberships/${this.currConvID}`, {username: user.username}, {
          headers: {
            Authorization: sessionStorage.getItem('ID')
          }
        });
        this.report = response.data.report;
        await this.fetchGroup();
        this.closeSearchBar();
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
        this.report = "";
      }, 2500)
    },

    async leaveGroup() {
      this.error = null;
      try {
        await this.$axios.delete(`/memberships/${this.currConvID}/members/${this.myID}`, {
          headers: {
            Authorization: sessionStorage.getItem('ID')
          }
        });
        this.$router.push({path: `/conversations`})
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
        if (this.newMember.length === 0) {
          this.searchResults = []
        }
        try {
          let response = await this.$axios.get(`/users/${this.myID}/search?username=${this.newMember}`, {
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
    }
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
    <ErrorMsg v-if="error" :msg="error" />
    <NotificationMsg v-if="report" :message="report" />

    <div class="w-75 h-auto align-items-center">
      <div class="text-white">
        <div class="position-absolute" style="top: 75px; left: 15px; z-index: 10;">
          <router-link :to="'/conversations/' + currConvID" class="btn btn-link text-white">
            <i class="fas fa-arrow-left fa-3x" />
          </router-link>
        </div>
      </div>

      <div
        class="text-center position-absolute d-flex flex-column p-3 rounded-3 gap-3"
        style="top: 10%; bottom: 10%; width: 30%; height: 80%; left: 35%; right: 35%; background-color: white; opacity: 0.9"
      >
        <div class="d-flex align-items-center justify-content-between w-100">
          <div class="d-flex align-items-center">
            <img
              :src="currGroup.photo || 'https://developer.jboss.org/images/jive-sgroup-default-portrait-large.png'"
              alt="Profile pic" class="profile-pic-header me-3"
            >
            <strong style="font-size: large; color: black"> {{ currGroup.name }} </strong>
          </div>
          <div>
            <i class="fas fa-user-plus text-primary me-3" style="cursor: pointer" @click="openSearchBar" />
            <i class="fas fa-sign-out-alt text-danger" style="cursor: pointer" @click="leaveGroup" />
          </div>
        </div>


        <div class="member-list w-100">
          <ul>
            <li v-for="member in members" :key="member.id" class="d-flex align-items-center p-2">
              <img :src="member.photo || 'https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg'" alt="Profile pic" class="rounded-circle me-2" width="40" height="40">
              <strong class="text-black"> {{ member.username }} </strong>
            </li>
          </ul>
        </div>

        <div class="w-100 mt-auto">
          <button class="rounded-3 w-100 mb-1" @click="openPhotoBar">Change group image</button>
          <button class="rounded-3 w-100 mb-2" @click="openNameBar">Change group name</button>
        </div>
      </div>
    </div>

    <div v-if="showPhotoBar" class="overlay">
      <div class="photo-box position-relative">
        <h3>Upload Profile Photo</h3>
        <div v-if="newPhoto" class="image-preview">
          <img :src="newPhoto" alt="Preview">
        </div>
        <input type="file" accept="image/*" @change="onFileChange">
        <button :disabled="!selectedFile" @click="setGroupPhoto">Upload</button>
        <button @click="closePhotoBar">Cancel</button>
      </div>
    </div>

    <div v-if="showUserSearch" class="overlay">
      <div class="search-box position-relative">
        <input v-model="newMember" placeholder="Search for a user to add..." @input="searchUsers">
        <ul>
          <li v-for="user in searchResults" :key="user.id" class="mb-0" @click="addToGroup(user)">
            {{ user.username }}
          </li>
        </ul>
        <button @click="closeSearchBar">Cancel</button>
      </div>
    </div>

    <div v-if="showNameBar" class="overlay">
      <div class="groupName-box position-relative">
        <input v-model="newName" placeholder="Enter a new group name..." @keyup.enter="setGroupName">
        <button @click="closeNameBar">Cancel</button>
      </div>
    </div>
  </div>
</template>

<style>

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

.groupName-box {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 30%;
  text-align: center;
}

.groupName-box input {
  width: 100%;
  padding: 10px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

ul {
  padding-inline-start: 0;
}

.member-list {
  height: 70vh;
  overflow-y: auto;
  border-radius: 10px;
  border: 1px solid #ddd;
  background-color: #f9f9f9;
  width: auto;
  color: black;
}

.member-list li {
  display: block;
  height: 5em;
  text-align: left;
  border: none;
  background-color: white;
  cursor: pointer;
  border-bottom: 1px solid #ddd;
  text-decoration: none;
  color: black;
  font-family: inherit;
}

.photo-box {
  background: white;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  color: black;
}

.image-preview img {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 10px;
}
</style>