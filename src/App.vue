<template>
  <b-container id="app">
    <b-row align-h="center">
      <b-col cols="1">
        <img
          class="text-right"
          src="../public/greek-column.svg"
          alt="greek pillar"
          height="50"
          width="50"
        >
      </b-col>
      <b-col>
        <h1>{{ $t('title') }}</h1>
      </b-col>
    </b-row>
    <router-view />
  </b-container>
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
  },
  created: function () {
    console.log(process.env)
  }
}
</script>
