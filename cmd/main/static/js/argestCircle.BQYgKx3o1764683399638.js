var i="varying vec2 vUv;\r\nvoid main(){\r\n\tvUv=uv;\r\n\tgl_Position=projectionMatrix*modelViewMatrix*vec4(position,1.);\r\n}";export{i as a};
