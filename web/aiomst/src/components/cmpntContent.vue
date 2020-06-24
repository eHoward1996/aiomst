<template>
  <v-content>
    <v-container fluid v-if="page == 'album'">
      <v-row>
        this this this
      </v-row>
    </v-container>
    <!-- <v-container fluid v-if="page == 'artist'">
      <v-row>
        
      </v-row>
    </v-container>
    <v-container fluid v-if="page == 'song'">
      <v-row>
      </v-row>
    </v-container> -->
  </v-content>
</template>

<script>
export default {
  name: 'Content',
  components: {},
  props: ['page'],
  mounted: function(page) {
    switch (page) {
      case 'album':

        break;
      case 'artist':
        break;
      case 'song':
        break;
    }
  },
  created() {
    this.$vuetify.theme.dark = true
  },
}
</script>