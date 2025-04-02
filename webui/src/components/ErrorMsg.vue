<script>
export default {
	props: ['msg']
}
</script>

<template>
  <div class="alert alert-danger" role="alert">
    {{ msg }}
  </div>
</template>

<style>

.alert {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgb(230, 0, 0);
  color: white;
  padding: 10px 20px;
  border-radius: 5px;
  font-size: 16px;
  font-weight: bold;
  z-index: 2000;
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateX(-50%) translateY(-10px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}
</style>
