<template>
  <b-col
    v-if="type === 'card'"
    cols="12"
    md="3"
    class="mb-3"
  >
    <b-card
      :title="book.title"
    >
      <b-link
        :to="{ name: 'book', params: { book: idSlug } }"
      >
        <b-card-img-lazy
          :src="coverImageLg"
        />
      </b-link>
      <b-card-text>
        <p
          v-if="published"
        >
          {{ $t('published') }} {{ published }}
        </p>
        <p
          v-if="authors"
        >
          {{ $t('attribution') }} {{ authors }}
        </p>
      </b-card-text>
      <b-button
        @click="read"
      >
        {{ $t('read.imperative') }}
      </b-button>
    </b-card>
  </b-col>

  <b-list-group-item
    v-else-if="type === 'list'"
    :to="{ name: 'book', params: { book: idSlug } }"
  >
    <h3>
      {{ book.title }}
      <small
        v-if="authors"
        class="text-muted"
      >
        {{ $t('attribution') }} {{ authors }}
      </small>
    </h3>
    <p
      v-if="book.timestamp"
      class="text-muted"
    >
      {{ $t('read.pastTense') }} {{ $d(new Date(), 'short') }}
    </p>
  </b-list-group-item>

  <div v-else>
    <b-row
      class="mb-3"
    >
      <b-col
        xl="4"
        lg="5"
        md="6"
        sm="8"
        xs="12"
        class="text-center"
      >
        <b-img
          :src="coverImageLg"
          responsive
        />
      </b-col>
      <b-col>
        <h1> {{ book.title }}</h1>
        <p
          v-if="authors"
          class="text-muted"
        >
          {{ $t('attribution') }} {{ authors }}
        </p>
        <p
          v-if="published"
        >
          {{ $t('published') }} {{ published }}
        </p>
        <b-button
          @click="read"
        >
          {{ $t('read.imperative') }}
        </b-button>
      </b-col>
    </b-row>

    <b-row>
      <b-col>
        <blockquote class="blockquote">
          {{ book.description }}
        </blockquote>
      </b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  name: 'Book',
  i18n: {
    messages: {
      en: {
        attribution: 'By',
        published: 'Published',
        read: {
          imperative: "I've Read This",
          pastTense: 'Read'
        },
        review: 'Write A Review',
        error: 'Error',
        readSuccess: {
          title: 'Book Read',
          message: '{title} has been added to your feed'
        }
      }
    }
  },
  props: {
    book: {
      type: Object,
      required: true
    },
    axios: {
      type: Function,
      default: null // not always required, as you can't always read
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
    },

    idSlug () {
      return this.book.id.split('/').pop()
    },

    published () {
      // don't have a published date
      if (!this.book.published) {
        return null
      }
      // published date just so happent to be the unix epoch? fishy.
      const timestamp = new Date(this.book.published)
      if (timestamp.getTime() === 0) {
        return null
      }
      return timestamp.getFullYear()
    }
  },
  methods: {
    read () {
      if (this.axios === null) {
        this.errorToast(new TypeError('cannot perform http request'))
        return
      }

      const self = this
      const id = this.book.id.split('/')[2]
      this.lastRead = this.book
      this.axios.post('/book/' + id + '/read')
        .then(self.successToast)
        .catch(self.errorToast)
    },

    successToast () {
      this.$bvToast.toast(this.$t('readSuccess.message', { title: this.book.title }), {
        title: this.$t('readSuccess.title'),
        solid: true,
        variant: 'info',
        autoHideDelay: 5000,
        appendToast: true
      })
    },

    errorToast (error) {
      this.$bvToast.toast(error.message, {
        title: this.$t('error'),
        solid: true,
        variant: 'danger',
        autoHideDelay: 5000,
        appendToast: true
      })
    }
  }
}
</script>
