<template>
  <b-img
    :fluid="fluid"
    :width="width"
    :height="height"
    alt=""
    :src="icon"
    :rounded="rounded"
    :left="left"
    :right="right"
    :center="center"
  />
</template>

<script>
export default {
  name: 'UserIcon',
  props: {
    user: {
      type: Object,
      required: true
    },
    round: {
      type: Boolean,
      default: false
    },
    width: {
      type: String,
      default: ''
    },
    height: {
      type: String,
      default: ''
    },
    left: {
      type: Boolean,
      default: false
    },
    right: {
      type: Boolean,
      default: false
    },
    center: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    icon () {
      return (this.user && this.user.icon) || require('../../public/default-profile.jpg')
    },
    rounded () {
      return this.round ? 'circle' : '0'
    },
    fluid () {
      return this.width === '' && this.height === ''
    }
  }
}
</script>
