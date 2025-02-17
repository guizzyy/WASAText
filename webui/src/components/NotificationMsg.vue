<script>
export default {
  props: ['message'],
}
</script>

<template>
  <div class="notification-msg" role="alert">
    {{ message }}
  </div>
</template>

<style>
</style>