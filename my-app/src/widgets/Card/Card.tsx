import React, { useState } from 'react';
import Button from 'react-bootstrap/Button';
import CardBootstrap from 'react-bootstrap/Card';
import { Link } from 'react-router-dom';

interface CardProps {
    id: number;
    name: string;
    image: string;
    price: number;
}

const Card: React.FC<CardProps> = (props) => {
    const [isHovered, setIsHovered] = useState(false);

    return (
        <CardBootstrap style={{ width: '18rem', marginTop: '3rem', margin: '10% 10% 5% 10%', backgroundColor: '#6884c0', color: 'white', position: 'relative' }}>
            <CardBootstrap.Img variant="top" src={props.image} />

            <CardBootstrap.Body style={{ textAlign: 'center' }}>
                <CardBootstrap.Title style={{ marginBottom: '3rem' }}>{props.name}</CardBootstrap.Title>

                {/* Цена слева внизу */}
                <div style={{ position: 'absolute', bottom: '0', left: '0', padding: '1rem' }}>
                    {props.price} рублей
                </div>
            </CardBootstrap.Body>

            {/* Кнопка справа внизу */}
            <div style={{ position: 'absolute', bottom: '0', right: '0', display: 'flex', justifyContent: 'flex-end', padding: '0.5rem' }}>
            <Link to={`/roads/${props.id}`} style={{ marginRight: '10px' }}>
                <Button
                    style={{
                        backgroundColor: isHovered ? '#f0f0f0' : '#000000',
                        borderColor: isHovered ? '#151d2a' : '#A0AECD',
                        color: isHovered ? '#000000' : '#f0f0f0'
                    }}
                    onMouseOver={() => setIsHovered(true)}
                    onMouseOut={() => setIsHovered(false)}
                >
                    Подробнее
                </Button>
            </Link>
            </div>
        </CardBootstrap>
    );
}

export default Card;
