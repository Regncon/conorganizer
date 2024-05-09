import { Button, Paper, Typography } from '@mui/material';
import Link from 'next/link';

type Props = {};

const List = ({}: Props) => {
	return (
		<Paper>
			<Typography variant="h1">Mine arrengementer</Typography>

			<Button component={Link} href="/event/create" variant="outlined">
				Lag nytt arrangement
			</Button>
		</Paper>
	);
};

export default List;
