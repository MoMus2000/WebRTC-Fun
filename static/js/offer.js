let peerConnection = new RTCPeerConnection()

const socket = new WebSocket("ws://localhost:3000/offer/ws")
socket.addEventListener('open', function(event){
    console.log("Connection Establised")
})
socket.addEventListener('message', async function(event){
    await startConnection(event.data)
})

peerConnection.ontrack = (event) => {
    event.streams[0].getTracks().forEach((track) => {
    const el = document.getElementById('client')
    // remoteStream.addTrack(track);
    el.srcObject = event.streams[0]
    el.autoplay = true
    el.controls = true
    });
};

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
        socket.send(encryptedOffer)
    }
    catch (error){
        console.log("Error accessing the webcam", error)
    }
}

async function startConnection(encryptedAnswer){
    // let encryptedAnswer = document.getElementById("answerForm").value
    console.log(JSON.parse(atob(encryptedAnswer)))
    await peerConnection.setRemoteDescription(JSON.parse(atob(encryptedAnswer)))
    console.log("Connection successful")
    // peerConnection.ontrack = function (event) {
    //     const el = document.getElementById('client')
    //     console.log(event)
    //     el.srcObject = event.streams[0]
    //     el.autoplay = true
    //     el.controls = true
    // }
}
