
let peerConnection = new RTCPeerConnection()

peerConnection.ontrack = function (event) {
        const el = document.getElementById('client')
        console.log(event)
        el.srcObject = event.streams[0]
        el.autoplay = true
        el.controls = true
    }

async function startWebcam(){
    try{
        const stream = await navigator.mediaDevices.getUserMedia({video : true})
        const videoElement = document.getElementById("webcam")
        const clientElement = document.getElementById("client")
        videoElement.srcObject = stream;

        let tracks = stream.getTracks();
        for(i=0; i<tracks.length; i++){
            peerConnection.addTrack(tracks[i], stream)
        }

        const offer = await peerConnection.createOffer()

        await peerConnection.setLocalDescription(offer);
        
        const encryptedOffer = btoa(JSON.stringify(peerConnection.localDescription))

        // const encryptedOffer = btoa(String.fromCharCode(...new Uint8Array(await crypto.subtle.encrypt({ name: 'AES-GCM', iv: new Uint8Array(12) }, await crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, ['encrypt', 'decrypt']), new TextEncoder().encode(JSON.stringify(offer))))));

        const freeform = document.getElementById("freeform")
        freeform.innerHTML = encryptedOffer
    }
    catch (error){
        console.log("Error accessing the webcam", error)
    }
}

async function startConnection(){
    let encryptedAnswer = document.getElementById("answerForm").value
    // console.log(JSON.parse(atob(encryptedAnswer)))
    // await peerConnection.setRemoteDescription(JSON.parse(atob(encryptedAnswer)))
    // console.log("Connection successful")
    // peerConnection.ontrack = function (event) {
    //     const el = document.getElementById('client')
    //     console.log(event)
    //     el.srcObject = event.streams[0]
    //     el.autoplay = true
    // //     el.controls = true
    // }
    peerConnection.onicecandidate = async (event) => {
    //Event that fires off when a new answer ICE candidate is created
        if(event.candidate){
            console.log('Adding answer candidate...:', event.candidate)
            document.getElementById('answerForm2').value = btoa(JSON.stringify(peerConnection.localDescription))
        }
    };

    await peerConnection.setRemoteDescription(JSON.parse(atob(encryptedAnswer)));

    let answer = await peerConnection.createAnswer();
    await peerConnection.setLocalDescription(answer); 
}
