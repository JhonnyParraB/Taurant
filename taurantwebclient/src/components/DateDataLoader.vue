<template>
  <v-row align="center" class="list px-3 mx-auto">
    <v-form v-model="valid">
        <div id="getUpdated">
            <h3>Get updated with the last information of your restaurant!</h3>
        </div>
        <div id="datePicker">
            <v-row justify="center">
                <v-date-picker color="#4b0082" v-model="picker"></v-date-picker>
            </v-row>
        </div>
        <div id="emailField">
        <v-text-field
            v-model="email"
            :rules="emailRules"
            label="E-mail"
            required
          ></v-text-field>
        </div>
        <div id="getUpdatedBtn">
            <v-btn
            class="mr-4"
            :disabled="!valid"
            outlined
            color="indigo"
            @click="loadDateData()"
            >
                Get Updated
            </v-btn>
        </div>
    </v-form>
  </v-row>
</template>


<script>
export default {
  name: "date-data-loader",
  data() {
    return {
      valid: false,
      email: '',
      picker: new Date().toISOString().substr(0, 10),
      emailRules: [
        v => !!v || 'E-mail is required',
        v => /.+@.+/.test(v) || 'E-mail must be valid',
      ],
    };
  },
  methods: {
    loadDateData(){
        
    },
  },
  mounted() {
    
  }
};
</script>

<style>
#getUpdated{
    margin-bottom: 5vh;
}
#emailField{
    margin-bottom: 5vh;
}
#datePicker{
    margin-bottom: 5vh;
}
#getUpdatedBtn{
    text-align: center;
}
</style>