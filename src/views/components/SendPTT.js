import FormRecipient from "./generic/FormRecipient.js";

export default {
    name: 'send-ptt',
    components: {
        FormRecipient
    },
    data() {
        return {
            phone: '',
            type: window.TYPEUSER,
            loading: false,
        }
    },
    computed: {
        phone_id() {
            return this.phone + this.type;
        }
    },
    methods: {
        openModal() {
            $('#modalPTTSend').modal({
                onApprove: function () {
                    return false;
                }
            }).modal('show');
        },
        async handleSubmit() {
            try {
                let response = await this.submitApi()
                showSuccessInfo(response)
                $('#modalPTTSend').modal('hide');
            } catch (err) {
                showErrorInfo(err)
            }
        },
        async submitApi() {
            this.loading = true;
            try {
                let payload = new FormData();
                payload.append("phone", this.phone_id)
                payload.append("audio", $("#file_ptt")[0].files[0])
                const response = await window.http.post(`/send/ptt`, payload)
                this.handleReset();
                return response.data.message;
            } catch (error) {
                if (error.response) {
                    throw new Error(error.response.data.message);
                }
                throw new Error(error.message);
            } finally {
                this.loading = false;
            }
        },
        handleReset() {
            this.phone = '';
            this.type = window.TYPEUSER;
            $("#file_ptt").val('');
        },
    },
    template: `
    <div class="blue card" @click="openModal()" style="cursor: pointer">
        <div class="content">
            <a class="ui blue right ribbon label">Send</a>
            <div class="header">Send Voice Message (PTT)</div>
            <div class="description">
                Send voice message (Push-to-Talk) to user or group
            </div>
        </div>
    </div>
    
    <!--  Modal SendPTT  -->
    <div class="ui small modal" id="modalPTTSend">
        <i class="close icon"></i>
        <div class="header">
            Send Voice Message (PTT)
        </div>
        <div class="content">
            <form class="ui form">
                <FormRecipient v-model:type="type" v-model:phone="phone"/>
                <div class="field" style="padding-bottom: 30px">
                    <label>Voice Message</label>
                    <input type="file" style="display: none" accept="audio/*" id="file_ptt"/>
                    <label for="file_ptt" class="ui positive medium green left floated button" style="color: white">
                        <i class="ui upload icon"></i>
                        Upload 
                    </label>
                </div>
            </form>
        </div>
        <div class="actions">
            <div class="ui approve positive right labeled icon button" :class="{'loading': this.loading}"
                 @click="handleSubmit">
                Send
                <i class="microphone icon"></i>
            </div>
        </div>
    </div>
    `
}
