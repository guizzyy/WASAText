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
      userWanted: "",
      convs: [
        {
          conversation: {
            ID : 0,
            type: "",
            name: "",
            photo: "",
            msg_unread: 0,
            last_mess: "",
            date_time: ""
          }
        }
      ],

      showLoading: false,
      showSearchInput: false,
    }
  },

  mounted() {
    this.getConversations();
  },

  methods: {
    async doLogout() {
      sessionStorage.clear();
      this.$router.push({ path: "/" });
    },

    async getConversations() {
      this.error = null;
      this.showLoading = true;
      try {
        let response = await this.$axios.get("/conversations", {
          headers: {
            Authorization: sessionStorage.getItem("ID"),
          }
        });
        this.convs = response.data;
        this.showLoading = false;
      } catch (e) {
        this.showLoading = false;
        if (e.response && e.response.status === 400) {
          this.error = "Failed to get conversations.";
        } else if (e.response && e.response.status === 404) {
          this.error = "User Not Found";
        } else if (e.response && e.response.status === 500) {
          this.error = "Server Error, please try again";
        } else {
          this.error = e.toString();
        }
        setTimeout(() => {
          this.error = null;
        }, 3000)
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

    async searchUsers() {
    },

    async setMyUsername(){},

    async toggleShowSearchInput() {
      this.showSearchInput = !this.showSearchInput;
      if (!this.showSearchInput) {
        this.userWanted = ""
      }
    }

  }

}
</script>

<template>

</template>

<style>
</style>
