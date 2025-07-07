package elements

import "html/template"

func Head() template.HTML {
	return template.HTML(`<head>
	<meta charset="UTF-8" name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="icon" href="../src/favicon.svg" />
	<title>universe</title>
	<style>
		 * {
			  box-sizing: border-box;
			  margin: 0;
			  scrollbar-width: none;
			  -ms-overflow-style: none;
			  user-select: none;
			  -webkit-user-select: none;
			  -moz-user-select: none;
			  -ms-user-select: none;
		 }
		 *::-webkit-scrollbar {
			  display: none;
		 }
		 html {
			  scroll-behavior: smooth;
		 }
		 body {
			  color: #ffffff;
			  background: #000000;
			  overflow: hidden;
			  height: 100vh;
			  width: 100vw;
			  margin: 0;
			  padding: 0;
			  border: medium solid #0000ff;
			  border-radius: 0.3125em;
			  font-family: 'Roboto', sans-serif;
           display: flex;
			  flex-direction: column;
		 }
	</style>
</head>`)
}
