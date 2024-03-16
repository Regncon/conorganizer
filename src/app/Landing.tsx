import * as React from 'react';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

export default function Landing() {
	return (
		<Box
			sx={{
				position: 'relative', // Set the Box to have relative positioning
				display: 'flex',
				justifyContent: 'center',
				alignItems: 'center',
				height: '100vh',
				bgcolor: '#eee943',
			}}
		>
			<Typography variant="h1" sx={{ position: 'relative', zIndex: 1 }}>
				BANAN
			</Typography>
			<Box
				component="img"
				src="/banan.png"
				alt="Banan"
				sx={{
					position: 'absolute',
					top: '45%', // Position the image right below the text
					left: '50%', // Center the image horizontally
					transform: 'translateX(-50%)', // Ensure it's centered by adjusting for its own width
					zIndex: 2,
					width: 1300, // Adjust width as needed
					height: 'auto', // Adjust height as needed, 'auto' keeps aspect ratio
				}}
			/>
		</Box>
	);
}
