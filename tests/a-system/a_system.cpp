#include <cstdio>

#include "a_system_messages.h"

namespace
{

const std::runtime_error no_error("");

a_system::Timestamp now()
{
    return a_system::Timestamp{0, 0};
}

} // anon

class service_impl : public a_system::Service
{
    std::string _key;
    a_system::ImageID _id = 0;

    void check(const std::string &tok)
    {
        if (tok != _key)
        {
            throw std::runtime_error("invalid service key");
        }
    }

public:
    service_impl(const std::string &service_key)
        : _key(service_key)
    {
    }

    ~service_impl() override
    {
    }

    a_system::Service_Auth_Result Auth(const std::string&, const a_system::Timestamp&) override
    {
        return a_system::Service_Auth_Result{_key, no_error};
    }

    a_system::Service_GetServerTime_Result GetServerTime(const std::string& key) override
    {
        check(key);
        return a_system::Service_GetServerTime_Result{now(), no_error};
    }

    a_system::Service_Hello_Result Hello(const std::string& name) override
    {
        return a_system::Service_Hello_Result{"Good to go "+name, no_error};
    }

    void Noop() override
    {
    }

    a_system::Service_PostImage_Result PostImage(const std::string& key, const a_system::Image& image) override
    {
        check(key);
        return a_system::Service_PostImage_Result{++_id, no_error};
    }

    a_system::Service_GetImage_Result GetImage(const std::string& key, const a_system::ImageID& id) override
    {
        check(key);
        if (id == 0) {
            const std::runtime_error error("image not found");
            return a_system::Service_GetImage_Result{a_system::Image{}, error};
        } else {
            return a_system::Service_GetImage_Result{a_system::Image{}, no_error};
        }
    }

    void Reset() override
    {
        _key = "";
    }
};

void run(a_system::Service &service)
{
}

int main()
{
    {
        auto m = a_system::message::Service_Reset{};
        printf("0x%x\n", m._header._msg);
    }

    service_impl service("123");
    run(service);
}
