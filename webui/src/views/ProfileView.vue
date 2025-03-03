<script>
import {RouterLink} from "vue-router";

export default {
  components: {
    data() {
      return {
        error: null,
        ID: sessionStorage.getItem("ID"),
        username: sessionStorage.getItem("username"),
        photo: sessionStorage.getItem("photo"),

        showLoading: false,
        showSearchBar: true,
      }
    }
  },
  methods: {
    async logout() {

    },

    async setMyUserName() {

    },

    async setMyPhoto() {

    },

    async createGroup() {

    },
  }
}
</script>

<template>

  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA Text</a>
    <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
  </header>

  <div class="text-center position-absolute d-flex flex-column justify-content-between align-items-center p-3"
       style="top: 10%; bottom: 10%; width: 30%; height: 80%; left: 35%; right: 35%; background-color: white; opacity: 0.9">

    <div class="profile-pic">
    </div>

    <div style="flex-grow: 1;">
      <span> {{ this.username }}</span>
    </div>



    <div class="w-100"> <!-- Ensures buttons take full width -->
      <button class="rounded-3 w-100 mb-1">Change your profile image</button>
      <button class="rounded-3 w-100 mb-2">Change your username</button>
      <button class="rounded-3 w-100 mb-3">Logout</button>
    </div>

  </div>

</template>

<style>

.profile-pic {
  margin-top: 50px;
  width: 250px;
  height: 250px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 20px;
  background-color: black;
}

button {
  color: white;
  background-color: #0b5ed7;
}

button:hover {
  background-color: #06357a;
}
</style>