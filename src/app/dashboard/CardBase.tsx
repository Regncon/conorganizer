'use client';
import Card from '@mui/material/Card';
import CardActionArea from '@mui/material/CardActionArea';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import type { Route } from 'next';
import CardActions from '@mui/material/CardActions';
import DeleteForeverOutlinedIcon from '@mui/icons-material/DeleteForeverOutlined';
import Box from '@mui/material/Box';
import { IconButton } from '@mui/material';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { collection, deleteDoc, doc } from 'firebase/firestore';
type Props = {
	href: Route;
	title: string;
	description: string;
	img: string;
	imgAlt: string;
	docId?: string;
};

const CardBase = ({ title, img, imgAlt, description, href, docId }: Props) => {
	const router = useRouter();
	const [disableRipple, setDisableRipple] = useState<boolean>(false);
	useEffect(() => {
		router.prefetch(href);
	});
	const handleActionClick = () => {
		router.push(href);
	};
	const handleDeleteClick = async () => {
		if (firebaseAuth.currentUser?.uid && docId) {
			const eventRef = doc(db, 'users', firebaseAuth.currentUser?.uid, 'my-events', docId);
			await deleteDoc(eventRef);
		}
	};
	return (
		<Card sx={{ maxWidth: 345 }}>
			<CardActionArea onClick={handleActionClick} disableRipple={disableRipple}>
				<CardMedia component="img" height={130} image={img} alt={imgAlt} />
				<CardContent>
					<Typography gutterBottom variant="h5" component="div">
						{title}
					</Typography>
					<Typography variant="body2" color="text.secondary">
						{description}
					</Typography>
				</CardContent>
				{docId ?
					<CardActions>
						<IconButton
							sx={{ placeSelf: 'end', color: '#f95e5e', padding: '1rem' }}
							onClick={(e) => {
								e.stopPropagation();
								setDisableRipple(true);
								handleDeleteClick();
							}}
						>
							<DeleteForeverOutlinedIcon />
						</IconButton>
					</CardActions>
				:	null}
			</CardActionArea>
		</Card>
	);
};

export default CardBase;
