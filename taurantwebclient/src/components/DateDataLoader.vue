<template>
  <v-row align="center" class="list px-3 mx-auto">
    <v-form v-model="valid">
        <div id="getUpdated">
            <h3>Get updated with the last information of your restaurant!</h3>
        </div>
        <div id="datePicker">
            <v-row justify="center">
                <v-date-picker color="#4b0082" v-model="date"></v-date-picker>
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

        <div class="text-center">
        <v-btn
            class="mr-4"
            :disabled="!valid"
            outlined
            color="indigo"
            @click="loadDateData()"
        >
          Get Updated
        </v-btn>

      <v-dialog
        v-model="dialogSuccesful"
        width="500"
      >   
      <v-card>
        <v-card-title class="headline primary">
          Request received
        </v-card-title>

        <v-card-text>
          Your request was received, we will send you a message to "{{this.email}}" with the result of the operation.
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="dialogSuccesful = false"
          >
            Got it
          </v-btn>
        </v-card-actions>
        </v-card>
        </v-dialog>



        <v-dialog
        v-model="dialogError"
        width="500"
      >   
      <v-card>
        <v-card-title class="headline error">
          Request Error
        </v-card-title>

        <v-card-text>
          We are responding to many requests. Try it again later. If the problem persists please contact the Taurant administrator.
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="error"
            text
            @click="dialogError = false"
          >
            Got it
          </v-btn>
        </v-card-actions>
        </v-card>
        </v-dialog>

        
        </div>
    </v-form>
  </v-row>
</template>


<script>
import DateDataLoaderService from "../services/DateDataLoaderService";
export default {
  name: "date-data-loader",
  data() {
    return {
      valid: false,
      dialogSuccesful: false,
      dialogError: false,
      email: '',
      date: new Date().toISOString().substr(0, 10),
      emailRules: [
        v => !!v || 'E-mail is required',
        v => /.+@.+/.test(v) || 'E-mail must be valid',
      ],
    };
  },
  methods: {
    loadDateData(){
        DateDataLoaderService.loadDayData(new Date(this.date).getTime() / 1000, this.email)
        .then(() => {
          this.dialogSuccesful = true;
        })
        .catch((e) => {
          console.log(e);
          this.dialogError = true;
        });
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
</style>