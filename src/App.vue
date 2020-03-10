<template>
  <div id="app">
    <div id="title">
      <h1>{{ $t('title') }}</h1>
      <img
        src="../public/greek-column.svg"
        alt="greek pillar"
        height="50"
        width="50"
      >
    </div>
    <div id="renderBox">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  i18n: {
    messages: {
      en: {
        title: 'exlibris'
      }
    }
  }
}
</script>

<style>
#app {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

h1 {
  text-align: center;
}

</style>
