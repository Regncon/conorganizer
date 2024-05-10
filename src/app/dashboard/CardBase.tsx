import Card from '@mui/material/Card';
import CardActionArea from '@mui/material/CardActionArea';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import Link from 'next/link';

type Props = {
	href: string;
	title: string;
	description: string;
	img: string;
	imgAlt: string;
};

const CardBase = ({ title, img, imgAlt, description, href }: Props) => {
	return (
		<Card sx={{ maxWidth: 345, textDecoration: 'none' }} component={Link} href={href}>
			<CardActionArea>
				<CardMedia component="img" height={130} image={img} alt={imgAlt} />
				<CardContent>
					<Typography gutterBottom variant="h5" component="div">
						{title}
					</Typography>
					<Typography variant="body2" color="text.secondary">
						{description}
					</Typography>
				</CardContent>
			</CardActionArea>
		</Card>
	);
};

export default CardBase;
