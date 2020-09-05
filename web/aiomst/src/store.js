import axios from 'axios';
import Vue from 'vue';
import Vuex from 'vuex';
import {Howl} from 'howler';

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    gResp: null,
    playback: null,
    playlist: null,
    replayState: 0,
  },
  getters: {
    currentArtist: state => {
      return state.gResp["artists"].length === 1 ? state.gResp : null
    },      
    getArtists: state => {
      return state.gResp["artists"];
    },

    currentAlbum: state => {
      return state.gResp["albums"].length === 1 ? state.gResp : null
    },      
    getAlbums: state => {
      return state.gResp["albums"];
    },

    currentSong: state => {
      return state.playback
    },
    getSongs: state => {
      return state.gResp["songs"];
    },
    getPlaylist: state => {
      return state.playlist;
    }
  },
  mutations: {
    changeGoResp: (state, payload) => {
      state.gResp = payload.gResp;
    },
    setSongState: (state, payload) => {
      state.playback = payload.playback
    },
    setPlaylist: (state, payload) => {
      state.playlist = payload
    },
    updateReplayState: (state) => {
      state.replayState = state.replayState === 2 ? 0 : state.replayState + 1
    }
  },
  actions: {
    makeApiRequest: (context, navInfo) => {
      return new Promise((resolve, reject) => {
        axios
          .get(navInfo.path, {
            baseURL: "http://127.0.0.1:8090",
            params: navInfo.params,
          })
          .then((x) => {
            context.commit('changeGoResp', {
              gResp: x.data,
            });
            resolve();
          })
          .catch((x) => {
            console.log(x);
            reject();
          });
      });
    },
    streamAudio: (context, song) => {
      // Check if a request for the current playing song has been made
      // If so, ignore it
      var playing = context.getters.currentSong;
      if (playing && playing.song.id === song.id) {
        return
      }
      
      return new Promise((resolve, reject) => {
        axios
          .get("http://localhost:8090/stream?id=" + song.id, {
            headers: {
              "Accept": "text/event-stream",
            },
            responseType: "arraybuffer"
          })
          .then(response => {
            // If there is a song playing, stop it
            if (context.getters.currentSong) {
              var playback = context.getters.currentSong
              playback.howl.unload()
            }
                  
            var blob = new Blob(
              [response.data], 
              {type: response.headers["content-type"]}
            );

            var url = URL.createObjectURL(blob);
            var howlInfo = new Howl({
              src: [url],
              html5: true,
              volume: .5,
              preload: false,
              onend: () => {
                var list = context.getters.getPlaylist
                var indexCurr = list.indexOf(context.getters.currentSong.song)
                var next = list[0]

                switch (context.state.replayState) {
                  case 0:
                    if (indexCurr === list.length - 1) {
                      return
                    }
                    next = list[indexCurr + 1]
                    break
                  case 1:
                    var howlObj = context.state.playback.howl;
                    howlObj.stop()
                    howlObj.play(context.state.playback.howlId)
                    return
                  case 2:
                    if (indexCurr + 1 < list.length) {
                      next = list[indexCurr + 1]
                    }
                    break
                }
                context.dispatch('streamAudio', next)
              }
            });
      
            var howlId = howlInfo.play();            
            context.commit(
              'setSongState', {
                'playback': {
                  'howl': howlInfo,
                  'howlId': howlId,
                  'song': song,
                }
              }
            )
            resolve();
          })
          .catch(error => {
            console.error(error);
            reject();
          });
      }); 
    }
  }
});