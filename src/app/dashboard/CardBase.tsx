'use client';
import Card from '@mui/material/Card';
import CardActionArea from '@mui/material/CardActionArea';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import type { Route } from 'next';

type Props = {
	href: Route;
	title: string;
	description: string;
	img: string;
	imgAlt: string;
};

const CardBase = ({ title, img, imgAlt, description, href }: Props) => {
	const router = useRouter();
	useEffect(() => {
		router.prefetch(href);
	});
	const handleActionClick = () => {
		router.push(href);
	};
	return (
		<Card sx={{ maxWidth: 345 }}>
			<CardActionArea onClick={handleActionClick}>
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
