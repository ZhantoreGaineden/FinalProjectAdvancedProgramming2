import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  PawPrint,
  UserPlus,
  ShoppingCart,
  RefreshCcw,
  CheckCircle2,
  PlayCircle,
  Trash2,
  Search
} from "lucide-react";
import "./style.css";

const API = {
  pets: "/pet-api/api/pets",
  users: "/user-api/api/users",
  orders: "/order-api/api/orders",
};

function randomEmail() {
  return `zhantore${Date.now()}@example.com`;
}

function pretty(data) {
  return JSON.stringify(data, null, 2);
}

function getId(data) {
  return data?.id || data?.user_id || data?.pet_id || data?.order_id || data?.user?.id || data?.pet?.id || data?.order?.id || "";
}

function App() {
  const [message, setMessage] = useState("Frontend is ready. Click Run Full Demo.");
  const [loading, setLoading] = useState(false);

  const [pets, setPets] = useState([]);
  const [orders, setOrders] = useState([]);

  const [userId, setUserId] = useState("");
  const [petId, setPetId] = useState("");
  const [orderId, setOrderId] = useState("");

  const [userForm, setUserForm] = useState({
    full_name: "Zhantore Gaineden",
    email: randomEmail(),
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

  async function request(url, options = {}) {
    const response = await fetch(url, options);
    const text = await response.text();

    let data;
    try {
      data = text ? JSON.parse(text) : {};
    } catch {
      data = { raw: text };
    }

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${pretty(data)}`);
    }

    return data;
  }

  async function safeRun(title, fn) {
    try {
      setLoading(true);
      setMessage(`${title}...`);
      const data = await fn();
      setMessage(`${title} success:\n${pretty(data)}`);
      return data;
    } catch (err) {
      setMessage(`${title} failed:\n${err.message}`);
      return null;
    } finally {
      setLoading(false);
    }
  }

  async function registerUserRaw() {
    const email = userForm.email || randomEmail();

    const payload = {
      full_name: userForm.full_name,
      email,
      password: userForm.password,
    };

    const data = await request(`${API.users}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    const id = getId(data);
    if (id) setUserId(id);

    return data;
  }

  async function loginUserRaw() {
    return await request(`${API.users}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: userForm.email,
        password: userForm.password,
      }),
    });
  }

  async function createPetRaw() {
    const payload = {
      name: petForm.name,
      category: petForm.category,
      breed: petForm.breed,
      age: Number(petForm.age),
      price: Number(petForm.price),
      status: petForm.status,
    };

    const data = await request(API.pets, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    const id = getId(data);
    if (id) setPetId(id);

    await listPetsRaw();
    return data;
  }

  async function listPetsRaw() {
    const data = await request(API.pets);
    const list = Array.isArray(data) ? data : data?.pets || [];
    setPets(list);
    return data;
  }

  async function getPetRaw() {
    if (!petId) throw new Error("Pet ID is empty.");
    return await request(`${API.pets}/${petId}`);
  }

  async function updatePetRaw() {
    if (!petId) throw new Error("Pet ID is empty.");

    const payload = {
      name: petForm.name,
      category: petForm.category,
      breed: petForm.breed,
      age: Number(petForm.age),
      price: Number(petForm.price),
      status: petForm.status,
    };

    const data = await request(`${API.pets}/${petId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    await listPetsRaw();
    return data;
  }

  async function deletePetRaw() {
    if (!petId) throw new Error("Pet ID is empty.");

    const data = await request(`${API.pets}/${petId}`, {
      method: "DELETE",
    });

    setPetId("");
    await listPetsRaw();
    return data;
  }

  async function createOrderRaw() {
    if (!userId) throw new Error("User ID is empty. Register user first.");
    if (!petId) throw new Error("Pet ID is empty. Create pet first.");

    const data = await request(API.orders, {
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

    const id = getId(data);
    if (id) setOrderId(id);

    return data;
  }

  async function getOrderRaw() {
    if (!orderId) throw new Error("Order ID is empty.");
    return await request(`${API.orders}/${orderId}`);
  }

  async function listUserOrdersRaw() {
    if (!userId) throw new Error("User ID is empty.");

    const data = await request(`/order-api/api/users/${userId}/orders`);
    const list = Array.isArray(data) ? data : data?.orders || [];
    setOrders(list);
    return data;
  }

  async function updateOrderStatusRaw() {
    if (!orderId) throw new Error("Order ID is empty.");

    return await request(`${API.orders}/${orderId}/status`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        status: "paid",
        user_email: userForm.email,
      }),
    });
  }

  async function cancelOrderRaw() {
    if (!orderId) throw new Error("Order ID is empty.");

    return await request(`${API.orders}/${orderId}/cancel`, {
      method: "POST",
    });
  }

  async function runFullDemo() {
    try {
      setLoading(true);

      const email = randomEmail();
      setUserForm((prev) => ({ ...prev, email }));

      setMessage("Step 1/5: registering user...");
      const user = await request(`${API.users}/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          full_name: userForm.full_name,
          email,
          password: userForm.password,
        }),
      });

      const createdUserId = getId(user);
      if (!createdUserId) throw new Error(`User was created but ID was not found:\n${pretty(user)}`);
      setUserId(createdUserId);

      setMessage("Step 2/5: creating pet...");
      const pet = await request(API.pets, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: `Buddy ${Math.floor(Math.random() * 1000)}`,
          category: petForm.category,
          breed: petForm.breed,
          age: Number(petForm.age),
          price: Number(petForm.price),
          status: "available",
        }),
      });

      const createdPetId = getId(pet);
      if (!createdPetId) throw new Error(`Pet was created but ID was not found:\n${pretty(pet)}`);
      setPetId(createdPetId);

      setMessage("Step 3/5: creating order...");
      const order = await request(API.orders, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: createdUserId,
          user_email: email,
          items: [
            {
              pet_id: createdPetId,
              price: Number(petForm.price),
            },
          ],
        }),
      });

      const createdOrderId = getId(order);
      if (!createdOrderId) throw new Error(`Order was created but ID was not found:\n${pretty(order)}`);
      setOrderId(createdOrderId);

      setMessage("Step 4/5: updating order status...");
      const paidOrder = await request(`${API.orders}/${createdOrderId}/status`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          status: "paid",
          user_email: email,
        }),
      });

      setMessage("Step 5/5: loading lists...");
      const petList = await listPetsRaw();
      const userOrders = await request(`/order-api/api/users/${createdUserId}/orders`);
      setOrders(userOrders?.orders || []);

      setMessage(
        "FULL DEMO SUCCESS ✅\n\n" +
        `User ID: ${createdUserId}\n` +
        `Pet ID: ${createdPetId}\n` +
        `Order ID: ${createdOrderId}\n\n` +
        "Created user:\n" + pretty(user) + "\n\n" +
        "Created pet:\n" + pretty(pet) + "\n\n" +
        "Created order:\n" + pretty(order) + "\n\n" +
        "Paid order:\n" + pretty(paidOrder) + "\n\n" +
        "Pets list:\n" + pretty(petList)
      );
    } catch (err) {
      setMessage(`FULL DEMO FAILED ❌\n\n${err.message}`);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    safeRun("Load pets", listPetsRaw);
  }, []);

  return (
    <div className="page">
      <header className="hero">
        <div>
          <p className="badge">Advanced Programming 2 Bonus Frontend</p>
          <h1>Pet Store Frontend</h1>
          <p className="subtitle">
            Working React UI for User, Pet and Order REST API Gateways.
          </p>
        </div>

        <div className="heroActions">
          <button className="bigButton" disabled={loading} onClick={runFullDemo}>
            <PlayCircle size={20} />
            {loading ? "Running..." : "Run Full Demo"}
          </button>
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
          <div className="inline">
            <input
              value={userForm.email}
              onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
            />
            <button
              className="iconButton"
              onClick={() => setUserForm({ ...userForm, email: randomEmail() })}
              title="Generate new email"
            >
              <RefreshCcw size={16} />
            </button>
          </div>

          <label>Password</label>
          <input
            type="password"
            value={userForm.password}
            onChange={(e) => setUserForm({ ...userForm, password: e.target.value })}
          />

          <div className="buttons">
            <button disabled={loading} onClick={() => safeRun("Register user", registerUserRaw)}>
              Register
            </button>
            <button disabled={loading} className="secondary" onClick={() => safeRun("Login user", loginUserRaw)}>
              Login
            </button>
          </div>

          <label>User ID</label>
          <input
            placeholder="Will be filled after register"
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
          <select
            value={petForm.status}
            onChange={(e) => setPetForm({ ...petForm, status: e.target.value })}
          >
            <option value="available">available</option>
            <option value="sold">sold</option>
            <option value="reserved">reserved</option>
          </select>

          <div className="buttons">
            <button disabled={loading} onClick={() => safeRun("Create pet", createPetRaw)}>
              Create
            </button>
            <button disabled={loading} className="secondary" onClick={() => safeRun("List pets", listPetsRaw)}>
              List
            </button>
          </div>

          <div className="buttons">
            <button disabled={loading} className="secondary" onClick={() => safeRun("Get pet", getPetRaw)}>
              <Search size={16} />
              Get
            </button>
            <button disabled={loading} className="secondary" onClick={() => safeRun("Update pet", updatePetRaw)}>
              Update
            </button>
            <button disabled={loading} className="danger" onClick={() => safeRun("Delete pet", deletePetRaw)}>
              <Trash2 size={16} />
            </button>
          </div>

          <label>Pet ID</label>
          <input
            placeholder="Will be filled after create pet"
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
            Register user and create pet first. Then create order.
          </p>

          <button disabled={loading} onClick={() => safeRun("Create order", createOrderRaw)}>
            Create Order
          </button>

          <label>Order ID</label>
          <input
            placeholder="Will be filled after create order"
            value={orderId}
            onChange={(e) => setOrderId(e.target.value)}
          />

          <div className="buttons">
            <button disabled={loading} className="secondary" onClick={() => safeRun("Get order", getOrderRaw)}>
              Get
            </button>
            <button disabled={loading} className="secondary" onClick={() => safeRun("List user orders", listUserOrdersRaw)}>
              User Orders
            </button>
          </div>

          <button disabled={loading} className="success" onClick={() => safeRun("Mark order as paid", updateOrderStatusRaw)}>
            <CheckCircle2 size={16} />
            Mark as Paid
          </button>

          <button disabled={loading} className="danger full" onClick={() => safeRun("Cancel order", cancelOrderRaw)}>
            Cancel Order
          </button>
        </section>

        <section className="card wide">
          <div className="cardTitle">
            <PawPrint />
            <h2>Pets List</h2>
          </div>

          <div className="items">
            {pets.length === 0 ? (
              <p className="hint">No pets loaded yet.</p>
            ) : (
              pets.map((pet, index) => (
                <button
                  className="itemButton"
                  key={pet.id || index}
                  onClick={() => {
                    setPetId(pet.id || "");
                    setPetForm({
                      name: pet.name || "",
                      category: pet.category || "",
                      breed: pet.breed || "",
                      age: pet.age || 0,
                      price: pet.price || 0,
                      status: pet.status || "available",
                    });
                    setMessage(`Selected pet:\n${pretty(pet)}`);
                  }}
                >
                  <strong>{pet.name || "Unnamed pet"}</strong>
                  <span>{pet.category} • {pet.breed} • ${pet.price}</span>
                  <small>ID: {pet.id}</small>
                </button>
              ))
            )}
          </div>
        </section>

        <section className="card wide">
          <div className="cardTitle">
            <ShoppingCart />
            <h2>User Orders</h2>
          </div>

          <div className="items">
            {orders.length === 0 ? (
              <p className="hint">No orders loaded yet.</p>
            ) : (
              orders.map((order, index) => (
                <button
                  className="itemButton"
                  key={order.id || index}
                  onClick={() => {
                    setOrderId(order.id || "");
                    setMessage(`Selected order:\n${pretty(order)}`);
                  }}
                >
                  <strong>Order {order.id}</strong>
                  <span>Status: {order.status} • Total: ${order.total_price}</span>
                  <small>User ID: {order.user_id}</small>
                </button>
              ))
            )}
          </div>
        </section>

        <section className="card wide">
          <h2>API Response / Demo Result</h2>
          <pre>{message}</pre>
        </section>
      </main>
    </div>
  );
}

createRoot(document.getElementById("root")).render(<App />);
