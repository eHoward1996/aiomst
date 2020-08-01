<template>
  <v-app-bar app clipped-left>
    <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
    <v-toolbar-title class="align-center">
      <router-link to="/" style="text-decoration:none; color: white;">
        <span class="title">AIOMST</span>
      </router-link>
    </v-toolbar-title>
    <v-row align="center">
      <v-col cols="3"></v-col>
      <v-col cols="5">
        <v-text-field
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
          value=""
          @keyup.enter="handleReq">
        </v-text-field>
      </v-col>
      <v-col cols="2"></v-col>
    </v-row>
    <template v-slot:extension>
      <v-tabs align-with-title v-model="active_tab">
        <v-tab v-for="tab in tabs" :key="tab.id" @click="route(tab.name)">
          {{tab.text}}
        </v-tab>
      </v-tabs>
    </template>
  </v-app-bar>
</template>

<script>
export default {
  name: 'cmpntAppBar',
  data: () => ({
    query: "",
    tabs: [
      {id: 0, name: '/',   text: 'Home'},
      {id: 1, name: '/artist', text: 'Artists'},
      {id: 2, name: '/album',  text: 'Albums'},
      {id: 3, name: '/song',   text: 'Songs'},
    ]
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
    },
    route: function(p) {
      this.$router.push({
        path: p
      });
    }
  },
  created() {
    this.$vuetify.theme.dark = true
  },
}
</script>