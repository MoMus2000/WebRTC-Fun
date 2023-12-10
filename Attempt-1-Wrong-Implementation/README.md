# WebRTC-Fun
## Exploring WebRTC

### What is WebRTC and How Does It Work?
WebRTC is a peer-to-peer communication protocol with an API defined inside the browser for use with JavaScript.

WebRTC allows users to connect directly without an intermediary server. This can result in faster communication and cost savings, as server overhead can be avoided.

### My Understanding So Far:

For a video chat application:
1. You create an offer on your requesting client.
2. You send that offer to the accepting client through various means, such as WhatsApp, manual copy-pasting, or WebSockets.
3. The accepting client accepts the sent offer and generates an answer.
4. You send the answer back to the requesting client to establish a communication stream and complete the WebRTC handshake.

To enhance system robustness, other technologies are involved, such as STUN servers for tracking ICE candidates, STUN servers help in discovering the public IP and port of a device, especially when it's behind a NAT (Network Address Translation)., TURN servers if P2P is not possible, and signaling servers for exchanging offers and answers.
