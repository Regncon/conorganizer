import Typography from '@mui/material/Typography';
import AppBar from '@mui/material/AppBar';
import { Toolbar, IconButton, Button } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';

export default function navbar() {
	return (
		<AppBar sx={{ width: '100%' }} position="static">
			<Toolbar>
				<IconButton size="large" edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }}>
					<MenuIcon />
				</IconButton>
				<Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
					News
				</Typography>
				<Button color="inherit">Login</Button>
			</Toolbar>
		</AppBar>
	);
}
