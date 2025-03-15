<script>
import {RouterLink} from "vue-router";
import NotificationMsg from "../components/NotificationMsg.vue";

export default {
  components: {RouterLink, NotificationMsg},
  data: function() {
    return {
      error: null,
      ID: sessionStorage.getItem("ID"),
      username: sessionStorage.getItem("username"),
      photo: sessionStorage.getItem("photo") || "https://static.vecteezy.com/system/resources/previews/013/360/247/non_2x/default-avatar-photo-icon-social-media-profile-sign-symbol-vector.jpg",
      newUsername: "",
      newPhoto: "",
      notification: "",
      selectedFile: null,

      showLoading: false,
      showUsernameBar: false,
      showPhotoBar: false
    }
  },

  methods: {
    logout() {
      sessionStorage.clear();
      this.$router.push({path: '/'});
    },

    openUsernameBar() {
      this.showUsernameBar = true;
    },

    openPhotoBar() {
      this.showPhotoBar = true;
    },

    closeUsernameBar() {
      this.showUsernameBar = false;
      this.newUsername = "";
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

    async setMyUserName() {
      try {
        this.error = null;
        let response = await this.$axios.put(`users/${this.ID}/username`, {username: this.newUsername},{
          headers: { Authorization: sessionStorage.getItem("ID") }
            });
        this.notification = response.data.report;
        sessionStorage.setItem('username', this.newUsername);
        this.username = this.newUsername;
        this.closeUsernameBar();
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
        this.notification = "";
      }, 2500)
    },

    async setMyPhoto() {
      if (!this.selectedFile) {
        this.error = "Please upload a photo";
        return;
      }
      try {
        let formData = new FormData();
        formData.append('photo', this.selectedFile);
        let response = await this.$axios.put(`users/${this.ID}/photo`, formData, {
          headers : {
            "Content-type" : "multipart/form-data",
            Authorization : sessionStorage.getItem("ID")
          }
        });
        this.notification = response.data.message;
        console.log("url photo: ", response.data.photo)
        let newUrl = response.data.photo;
        sessionStorage.setItem("photo", newUrl);
        this.photo = newUrl;
        console.log(this.photo)
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
        this.notification = "";
      }, 2500)
    },
  },
}
</script>

<template>

  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-5">WASA Text</a>

    <div class="set-buttons d-flex align-items-center me-3">
      <button class="icon-btn" aria-label="Home">
        <router-link to="/conversations" class="icon-btn">
          Home
        </router-link>
      </button>
      <button class="icon-btn" aria-label="Profile">
        <router-link :to="'/users/' + this.ID" class="icon-btn">
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
    <NotificationMsg v-if="notification" :message="notification"></NotificationMsg>
  </header>


  <div class="w-75 h-auto align-items-center">
    <div class="text-center position-absolute d-flex flex-column justify-content-between align-items-center p-3 rounded-3"
         style="top: 10%; bottom: 10%; width: 30%; height: 80%; left: 35%; right: 35%; background-color: white; opacity: 0.9">
      <div>
        <img :src="photo" alt="Profile pic" class="profile-pic">
      </div>

      <div style="flex-grow: 1; color: black">
        <strong style="font-size: large"> {{this.username}} </strong>
      </div>


      <div class="w-100">
        <button class="rounded-3 w-100 mb-1" @click="openPhotoBar">Change your profile image</button>
        <button class="rounded-3 w-100 mb-2" @click="openUsernameBar">Change your username</button>
      </div>

      <div v-if="showUsernameBar" class="overlay">
        <div class="username-box position-relative">
          <input v-model="newUsername" placeholder="Enter a new username..." @keyup.enter="setMyUserName">
          <button @click="closeUsernameBar">Cancel</button>
        </div>
      </div>

      <div v-if="showPhotoBar" class="overlay">
        <div class="photo-box position-relative">
          <h3>Upload Profile Photo</h3>
          <div v-if="this.newPhoto" class="image-preview">
            <img :src="this.newPhoto" alt="Preview" />
          </div>
          <input type="file" @change="onFileChange" accept="image/*" />
          <button @click="setMyPhoto" :disabled="!selectedFile">Upload</button>
          <button @click="closePhotoBar">Cancel</button>
        </div>
      </div>
    </div>
  </div>

</template>

<style scoped>

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

.username-box {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 30%;
  text-align: center;
}

.username-box input {
  width: 100%;
  padding: 10px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
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

.profile-pic {
  width: 250px;
  height: 250px;
  border-radius: 50%;
  object-fit: cover;
  margin-top: 20px;
  margin-bottom: 20px;
  background-color: black;
}

button {
  margin: 5px;
  padding: 8px 15px;
  border: none;
  cursor: pointer;
}

button:disabled {
  background: gray;
  cursor: not-allowed;
}
</style>