--
-- Last.Backend LLC CONFIDENTIAL
-- __________________
--
-- [2014] - [2018] Last.Backend LLC
-- All Rights Reserved.
--
-- NOTICE:  All information contained herein is, and remains
-- the property of Last.Backend LLC and its suppliers,
-- if any.  The intellectual and technical concepts contained
-- herein are proprietary to Last.Backend LLC
-- and its suppliers and may be covered by Russian Federation and Foreign Patents,
-- patents in process, and are protected by trade secret or copyright law.
-- Dissemination of this information or reproduction of this material
-- is strictly forbidden unless prior written permission is obtained
-- from Last.Backend LLC.
--

local json = require("cjson")
local routes = lrucache:new(2048)

local lrucache = require("resty.lrucache")

local _M = {}

function get()
  return routes:get("localhost")
end

function set(data)
  for i = 1, #data do
    ngx.log(ngx.NOTICE,  "add endpoint: " .. tostring(data[i].endpoint) .. tostring(data[i].backends))
    routes:set(data[i].endpoint, "test")
  end

  ngx.log(ngx.NOTICE,  "added>")
  for k , v in pairs(data) do
    ngx.log(ngx.NOTICE, tostring(k).."  "..tostring(v))
  end
  ngx.log(ngx.NOTICE,  "added>")

  return true
end

function del(data)
  for i = 1, #data do
    ngx.log(ngx.NOTICE,  "del endpoint: " .. tostring(data[i].endpoint, data[i].backends))
    routes:del(data[i].endpoint)
  end
  return true
end

function parse()
  ngx.req.read_body()
  ngx.log(ngx.NOTICE,  "parse data: " .. tostring(ngx.req.get_body_data()))
  local ok, data = pcall(json.decode, ngx.req.get_body_data())
  if not ok then
    ngx.log(ngx.NOTICE,  "could not parse backends data: " .. tostring(ngx.req.get_body_data()))
    ngx.status = ngx.HTTP_BAD_REQUEST
    return
  end
  ngx.log(ngx.NOTICE, "parsed data: ".. tostring(#data))
  ngx.log(ngx.NOTICE,  "parsed>")
  for k , v in pairs(data) do
    ngx.log(ngx.NOTICE, tostring(k).."  "..tostring(v))
  end
  ngx.log(ngx.NOTICE,  "<parsed")
  return data
end

function _M.find()
  local host = ngx.var.host
  return routes:get(host)
end

function _M.call()

  if ngx.var.request_uri ~= "/" then
    ngx.status = ngx.HTTP_NOT_FOUND
    ngx.print("Only root allowed")
    return
  end

  if ngx.var.request_method == "GET" then
    ngx.status = ngx.HTTP_OK
    ngx.print(get())
    return
  end

  if ngx.var.request_method == "DEL" then
    -- Remove routes from cache
    local data = parse()
    if not data then
      return
    end
    del(data)
    ngx.status = ngx.HTTP_DELETED
  end

  if ngx.var.request_method == "POST" then

    local data = parse()
    if not data then
      return
    end

    local ok, err = set(data)
    if not ok then
      ngx.log(ngx.ERR, "can not save routes cache: " .. tostring(err))
      ngx.status = ngx.HTTP_BAD_REQUEST
      return
    end

    ngx.status = ngx.HTTP_CREATED
    return
  end

  ngx.status = ngx.HTTP_BAD_REQUEST
  ngx.print("Unsupported request type")
  return
end

return _M
