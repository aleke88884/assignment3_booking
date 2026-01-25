1. Project Proposal
Project Title

SmartBooking — Online Booking Management System

Project Relevance

Today, online booking systems are widely used in many domains such as hotel reservations, coworking spaces, meeting rooms, sport facilities, and service appointments. However, many existing solutions are either too complex, expensive, or not flexible enough for small and medium-sized businesses.

The goal of SmartBooking is to provide a simple, scalable, and user-friendly booking system that allows users to reserve resources online while giving administrators full control over availability, pricing, and bookings.

This project is relevant because it demonstrates real-world backend design, data modeling, concurrency handling, and modular system architecture using Go.

Target Users

End Users – people who want to book a resource (room, apartment, service, etc.)

Administrators – users who manage resources, schedules, and bookings

System Operators – responsible for maintaining the system and monitoring usage

Competitor Analysis

Existing booking platforms include:

Booking.com – powerful but overly complex and expensive

Airbnb – focused only on short-term accommodation

Calendly – suitable for time slots but not resource-based booking

SmartBooking aims to be:

simpler than large platforms

customizable for different booking scenarios

focused on backend architecture rather than UI complexity

Planned Features (High-Level)

User registration and authentication

Resource creation and management

Availability scheduling

Booking creation and cancellation

Conflict prevention (double booking protection)

Basic role-based access (user/admin)

UI workflows and final design decisions will be implemented in later milestones.

2. Architecture & Design
System Architecture

The system will follow a monolithic architecture, as required, with clearly separated internal modules.

Architecture Style:

Monolithic backend

REST API

Layered architecture (handlers → services → repositories)

Main Modules

Auth Module

User registration and login

Role management (user/admin)

User Module

User profile data

Booking history

Resource Module

Create and manage bookable resources

Define availability rules

Booking Module

Create bookings

Validate availability

Prevent overlapping reservations

Database Module

Data persistence

Transaction handling

Data Flow Overview

Client sends HTTP request

Handler validates request

Service applies business logic

Repository interacts with database

Response returned to client

Use-Case Diagram (Description)

Actors:

User

Admin

Use Cases:

Register / Login

View available resources

Create booking

Cancel booking

Manage resources (Admin)

View all bookings (Admin)

ER Diagram (Entities)

Entities:

User (id, name, email, role)

Resource (id, name, description, capacity)

Booking (id, user_id, resource_id, start_time, end_time, status)

Relationships:

User → Booking (one-to-many)

Resource → Booking (one-to-many)

UML Diagram (High-Level Classes)

UserService

BookingService

ResourceService

AuthService

UserRepository

BookingRepository

ResourceRepository

Each service communicates only with its corresponding repository.

3. Project Plan (Gantt – Weeks 7–10)
Week 7

Project requirements analysis

Architecture design

Database schema design

Git repository initialization

Week 8

Auth module skeleton

User module skeleton

Basic REST API setup

Initial ERD and UML diagrams

Week 9

Booking module logic

Availability validation

Resource management module

Basic testing

Week 10

Code refactoring

Documentation updates

Defense preparation

Task Distribution

Developer 1: Auth & User modules

Developer 2: Booking logic

Developer 3: Resource management & database layer
