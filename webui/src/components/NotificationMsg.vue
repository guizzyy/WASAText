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

.notification-msg {
  position: fixed; /* Always stays in the same place */
  top: 20px; /* Position near the top */
  left: 50%;
  transform: translateX(-50%); /* Centers it horizontally */
  background: rgb(0, 230, 0); /* Red background with slight transparency */
  color: white; /* White text for contrast */
  padding: 10px 20px;
  border-radius: 5px;
  font-size: 16px;
  font-weight: bold;
  z-index: 2000; /* Make sure it's above everything */
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateX(-50%) translateY(-10px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}

</style>