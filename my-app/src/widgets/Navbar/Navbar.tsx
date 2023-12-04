import { ChangeEvent } from 'react';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import Nav from 'react-bootstrap/Nav';
import { Navbar as NavB } from 'react-bootstrap';
import { useState } from 'react';

interface NavbarProps {
  onMinLengthChange?: (value: string) => void; // Define the prop type
}

const Navbar: React.FC<NavbarProps> = ({ onMinLengthChange }) => {
  const [minLength, setMinLength] = useState('');

  const handleMinLengthChange = (e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault()
    const value = e.target.value;
    setMinLength(value);
    if (onMinLengthChange !== undefined) {
      onMinLengthChange(value);
    }
  };
  const handleSearchSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    // Вызываем onMaxPriceChange при отправке формы
    if (onMinLengthChange && minLength.trim() !== '') {
      onMinLengthChange(minLength);
    }
  };

  return (
    <NavB expand="lg" className="bg-dark" style={{ backgroundColor: '#000000' }}>
      <Container fluid>
        <NavB.Brand href="#" className="text-light">Платные дороги</NavB.Brand>
        <NavB.Toggle aria-controls="navbarScroll" />
        <NavB.Collapse id="navbarScroll">
          <Nav className="me-auto my-2 my-lg-0" style={{ maxHeight: '100px' }} navbarScroll>
            <Nav.Link href="/" className="text-light">Главная</Nav.Link>
            <Nav.Link href="#action2" className="text-light">Корзина</Nav.Link>
          </Nav>
          <Form className="d-flex" id="search" onSubmit={handleSearchSubmit}>
            <Form.Control
              type="search"
              placeholder="Минимальная длина"
              className="me-2"
              aria-label="Search"
              value={minLength}
              onChange={handleMinLengthChange}
            />
             <Button
              variant="outline-light"
              onClick={(e) => {
                e.preventDefault();
                if (onMinLengthChange !== undefined) {
                  onMinLengthChange(minLength);
                }
              }} 
            >Поиск</Button>
          </Form>
        </NavB.Collapse>
      </Container>
    </NavB>
  );
}

export default Navbar;
