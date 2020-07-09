<template>
  <b-col
    v-if="type === 'card'"
    cols="12"
    md="3"
  >
    <b-card
      :title="book.title"
      :img-src="coverImageLg"
    >
      <b-card-text>
        {{ $t('attribution') }} {{ book.authors[0] }}
      </b-card-text>
      <b-button
        @click="readBook"
      >
        {{ $t('read') }}
      </b-button>
    </b-card>
  </b-col>
  <b-list-group-item
    v-else-if="type === 'list'"
  >
    <h3>{{ book.title }}</h3>
    <p
      v-if="authors"
      class="text-muted"
    >
      {{ $t('attribution') }} {{ authors }}
    </p>
  </b-list-group-item>
</template>

<script>
export default {
  name: 'Book',
  i18n: {
    messages: {
      en: {
        attribution: 'By',
        read: "I've Read This",
        review: 'Write A Review'
      }
    }
  },
  props: {
    book: {
      type: Object,
      required: true
    },
    type: {
      type: String,
      default: 'full'
    }
  },
  computed: {
    authors () {
      return this.book && this.book.authors && this.book.authors.join(', ')
    },

    coverImageLg: function () {
      return (this.book.covers && (this.book.covers.large || this.book.covers.medium || this.book.covers.small)) || require('../../public/default-cover.jpg')
    },

    coverImageSm: function () {
      return (this.book.covers && (this.book.covers.small || this.book.covers.medium || this.book.covers.large)) || require('../../public/default-cover.jpg')
    }
  },
  created () {
    console.log(this.book)
  },
  methods: {
    readBook: function () {
      this.$emit('read', this.book)
    }
  }
}
</script>
