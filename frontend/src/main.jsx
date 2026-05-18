import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { PawPrint, UserPlus, ShoppingCart, RefreshCcw, CheckCircle2 } from "lucide-react";
import "./style.css";

const API = {
  pets: "/pet-api/api/pets",
  users: "/user-api/api/users",
  orders: "/order-api/api/orders",
};

function App() {
  const [message, setMessage] = useState("Frontend is ready.");
  const [pets, setPets] = useState([]);
  const [userId, setUserId] = useState("");
  const [petId, setPetId] = useState("");
  const [orderId, setOrderId] = useState("");

  const [userForm, setUserForm] = useState({
    full_name: "Zhantore Gaineden",
    email: "zhantore@example.com",
    password: "123456",
  });

  const [petForm, setPetForm] = useState({
    name: "Buddy",
    category: "dog",
    breed: "Golden Retriever",
    age: 2,
    price: 500,
    status: "available",
  });

  const showResult = async (response) => {
    const text = await response.text();
    let data;

    try {
      data = JSON.parse(text);
    } catch {
      data = text;
    }

    if (!response.ok) {
      setMessage(`Error ${response.status}: ${JSON.stringify(data, null, 2)}`);
      return null;
    }

    setMessage(JSON.stringify(data, null, 2));
    return data;
  };

  const registerUser = async () => {
    const res = await fetch(`${API.users}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(userForm),
    });

    const data = await showResult(res);
    const id = data?.id || data?.user_id || data?.user?.id;

    if (id) {
      setUserId(id);
    }
  };

  const loginUser = async () => {
    const res = await fetch(`${API.users}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: userForm.email,
        password: userForm.password,
      }),
    });

    await showResult(res);
  };

  const createPet = async () => {
    const payload = {
      ...petForm,
      age: Number(petForm.age),
      price: Number(petForm.price),
    };

    const res = await fetch(API.pets, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    const data = await showResult(res);
    const id = data?.id || data?.pet_id || data?.pet?.id;

    if (id) {
      setPetId(id);
    }

    await listPets();
  };

  const listPets = async () => {
    const res = await fetch(API.pets);
    const data = await showResult(res);

    if (Array.isArray(data)) {
      setPets(data);
    } else if (Array.isArray(data?.pets)) {
      setPets(data.pets);
    }
  };

  const createOrder = async () => {
    if (!userId || !petId) {
      setMessage("Please enter or create userId and petId first.");
      return;
    }

    const res = await fetch(API.orders, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_id: userId,
        user_email: userForm.email,
        items: [
          {
            pet_id: petId,
            price: Number(petForm.price),
          },
        ],
      }),
    });

    const data = await showResult(res);
    const id = data?.id || data?.order_id || data?.order?.id;

    if (id) {
      setOrderId(id);
    }
  };

  const updateOrderStatus = async () => {
    if (!orderId) {
      setMessage("Please enter or create orderId first.");
      return;
    }

    const res = await fetch(`${API.orders}/${orderId}/status`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        status: "paid",
        user_email: userForm.email,
      }),
    });

    await showResult(res);
  };

  useEffect(() => {
    listPets();
  }, []);

  return (
    <div className="page">
      <header className="hero">
        <div>
          <p className="badge">Advanced Programming 2 Bonus Frontend</p>
          <h1>Pet Store Microservices</h1>
          <p className="subtitle">
            React frontend for Pet, User and Order API Gateways.
          </p>
        </div>
        <div className="heroIcon">
          <PawPrint size={56} />
        </div>
      </header>

      <main className="grid">
        <section className="card">
          <div className="cardTitle">
            <UserPlus />
            <h2>User</h2>
          </div>

          <label>Full name</label>
          <input
            value={userForm.full_name}
            onChange={(e) => setUserForm({ ...userForm, full_name: e.target.value })}
          />

          <label>Email</label>
          <input
            value={userForm.email}
            onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
          />

          <label>Password</label>
          <input
            type="password"
            value={userForm.password}
            onChange={(e) => setUserForm({ ...userForm, password: e.target.value })}
          />

          <div className="buttons">
            <button onClick={registerUser}>Register</button>
            <button className="secondary" onClick={loginUser}>Login</button>
          </div>

          <label>User ID</label>
          <input
            placeholder="Paste user id here"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
          />
        </section>

        <section className="card">
          <div className="cardTitle">
            <PawPrint />
            <h2>Pet</h2>
          </div>

          <label>Name</label>
          <input
            value={petForm.name}
            onChange={(e) => setPetForm({ ...petForm, name: e.target.value })}
          />

          <label>Category</label>
          <input
            value={petForm.category}
            onChange={(e) => setPetForm({ ...petForm, category: e.target.value })}
          />

          <label>Breed</label>
          <input
            value={petForm.breed}
            onChange={(e) => setPetForm({ ...petForm, breed: e.target.value })}
          />

          <div className="two">
            <div>
              <label>Age</label>
              <input
                type="number"
                value={petForm.age}
                onChange={(e) => setPetForm({ ...petForm, age: e.target.value })}
              />
            </div>

            <div>
              <label>Price</label>
              <input
                type="number"
                value={petForm.price}
                onChange={(e) => setPetForm({ ...petForm, price: e.target.value })}
              />
            </div>
          </div>

          <label>Status</label>
          <input
            value={petForm.status}
            onChange={(e) => setPetForm({ ...petForm, status: e.target.value })}
          />

          <div className="buttons">
            <button onClick={createPet}>Create Pet</button>
            <button className="secondary" onClick={listPets}>
              <RefreshCcw size={16} />
              List Pets
            </button>
          </div>

          <label>Pet ID</label>
          <input
            placeholder="Paste pet id here"
            value={petId}
            onChange={(e) => setPetId(e.target.value)}
          />
        </section>

        <section className="card">
          <div className="cardTitle">
            <ShoppingCart />
            <h2>Order</h2>
          </div>

          <p className="hint">
            Use created User ID and Pet ID to create an order.
          </p>

          <button onClick={createOrder}>Create Order</button>

          <label>Order ID</label>
          <input
            placeholder="Paste order id here"
            value={orderId}
            onChange={(e) => setOrderId(e.target.value)}
          />

          <button className="success" onClick={updateOrderStatus}>
            <CheckCircle2 size={16} />
            Mark Order as Paid
          </button>
        </section>

        <section className="card wide">
          <div className="cardTitle">
            <RefreshCcw />
            <h2>Pets List</h2>
          </div>

          <div className="pets">
            {pets.length === 0 ? (
              <p className="hint">No pets loaded yet.</p>
            ) : (
              pets.map((pet, index) => (
                <div className="petItem" key={pet.id || index}>
                  <strong>{pet.name || "Unnamed pet"}</strong>
                  <span>{pet.category} • {pet.breed} • ${pet.price}</span>
                  <small>ID: {pet.id}</small>
                </div>
              ))
            )}
          </div>
        </section>

        <section className="card wide">
          <h2>API Response</h2>
          <pre>{message}</pre>
        </section>
      </main>
    </div>
  );
}

createRoot(document.getElementById("root")).render(<App />);
