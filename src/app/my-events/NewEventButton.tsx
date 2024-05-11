'use client';
import IconButton from '@mui/material/IconButton';
import AddIcon from '@mui/icons-material/Add';
import { useRouter } from 'next/navigation';
import { createMyEventDoc } from './actions';
type Props = { docId: string };

const NewEventButton = ({ docId }: Props) => {
	const router = useRouter();
	router.prefetch(`/event/create/${docId}`);

	const handleClick = async () => {
		await createMyEventDoc(docId);
		router.push(`/event/create/${docId}`);
	};
	return (
		<IconButton
			sx={{
				border: '1px solid',
				borderColor: 'secondary.contrastText',
				position: 'absolute',
				zIndex: '22',
				right: '0.3125rem',
				top: '-0.375rem',
				'&:hover, &:focus, &': {
					backgroundColor: 'secondary.contrastText',
				},
			}}
			onClick={handleClick}
		>
			<AddIcon />
		</IconButton>
	);
};

export default NewEventButton;
