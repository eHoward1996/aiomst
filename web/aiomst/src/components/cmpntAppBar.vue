<template>
  <v-app-bar app clipped-left>
    <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
    <v-toolbar-title class="align-center">
      <span class="title">AIOMST</span>
    </v-toolbar-title>
    <v-row align="center">
      <v-col cols="3"></v-col>
      <v-col cols="5">
        <v-text-field
          append-icon=""
          color="#43A047"
          label="Search"
          prepend-inner-icon="mdi-magnify"
          hide-details="auto"
          dense
          outlined
          rounded
          single-line
          name="query"
          v-model="query"
          @keyup.enter="handleReq">
        </v-text-field>
      </v-col>
      <v-col cols="2"></v-col>
    </v-row>
  </v-app-bar>
</template>

<script>
// import axios from 'axios';

export default {
  name: 'AppBar',

  data: () => ({
    query: "",
  }),
  methods:  {
    handleReq: function ()  {
      if (this.query.length < 3 && this.query.length != 0) {
        this.gResp = "Min 3 Characters for Search."
        return
      }
      
      this.$router.push({
        path: 'search', 
        query: {q: this.query},
      })
    }
  },
  created() {
    this.$vuetify.theme.dark = true
  },
}
</script>