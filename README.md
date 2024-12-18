# Newsletter Microservice

This is a Newsletter Microservice that allows users to subscribe to topics of their interest, receive newsletters via email, and manage their preferences. The system consists of multiple components including subscription handling, email sending, and integration with open-source APIs for content fetching.

## Features

- **User Registration and Authentication**: Users can register and log in securely.
- **Topic Subscription**: Users can subscribe to different topics (e.g., Technology, Science, etc.).
- **Content Fetching**: The system fetches relevant content from open-source APIs based on user preferences.
- **Email Notifications**: Users receive weekly emails containing content fetched from the subscribed topics.
- **JWT Authentication**: Secure communication between clients and services using JWT tokens.
- **Environment Configuration**: Environment variables are used for flexible and secure configuration.

## Technologies Used

- **Go**: Backend service implemented using Go.
- **Gin**: Web framework for handling HTTP requests.
- **GORM**: ORM for database interaction.
- **PostgreSQL**: Database for storing user data and subscriptions.
- **JWT**: For authentication and token-based user validation.
- **SendGrid**: For sending subscription updates and newsletters.


## Setup

### Prerequisites

- Go 1.16 or later
- PostgreSQL database


