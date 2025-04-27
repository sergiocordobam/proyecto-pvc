FROM ruby:2.7
WORKDIR /app
COPY app.rb .
RUN gem install sinatra
CMD ["ruby", "app.rb", "-o", "0.0.0.0"]
