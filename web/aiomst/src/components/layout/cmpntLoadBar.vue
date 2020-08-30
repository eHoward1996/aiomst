<template>
  <v-container 
    fluid
    fill-height
    style="width: 50%; height: 25vh;">
    <v-row justify="center" align="center" class="mb-n12">
      <h2>LOADING...</h2>
    </v-row>
    <v-row align="center" class="mt-n12">
      <v-progress-linear
        indeterminate
        rounded
        color="green lighten-0">
      </v-progress-linear>
  </v-row>
  </v-container>
</template>

<script>
export default {
  name: 'cmpntLoadBar',
}
</script>