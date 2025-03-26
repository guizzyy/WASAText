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
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgb(0, 230, 0);
  color: white;
  padding: 10px 20px;
  border-radius: 5px;
  font-size: 16px;
  font-weight: bold;
  animation: fadeIn 0.3s ease-in-out;
  z-index: 9999;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateX(-50%) translateY(-10px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}

</style>