class Name:
    def __init__(self, first_name, middle_name, last_name):
        self.first = first_name
        self.middle = middle_name
        self.last = last_name

    def to_json(self):
        name = {
            "first" : self.first,
            "middle": self.middle,
            "last"  : self.last,
        }
        if len(name["middle"]) == 0:
            del name["middle"]
        return name


class Person(object):
    def __init__(self, internal_id=None, first_name=None, middle_name="", last_name=None, phone_number=None):
        self.id = internal_id
        if internal_id is None or len(internal_id) == 0:
            raise ValueError("INTERNAL_ID Must be set")
        if first_name is None or len(first_name) == 0:
            raise ValueError("FIRST_NAME Must be set")
        if last_name is None or len(last_name) == 0:
            raise ValueError("LAST_NAME Must be set")
        if phone_number is None or len(phone_number) == 0:
            # TODO: Add validator for Phone Number to ensure it's in format ###-###-####
            raise ValueError("PHONE_NUM Must be set")

        self.name = Name(first_name=first_name, middle_name=middle_name, last_name=last_name)
        self.phone = phone_number

    def to_json(self) -> object:
        """
        to_json will return an Object that can be written out to json with the
        json.dumps method. Remove's the middle name if it is not set.
        :return: object
        """
        return {
            "id"   : self.id,
            "name" : self.name.to_json(),
            "phone": self.phone
        }


class Error(object):
    def __init__(self, row, error):
        self.row = row
        self.message = str(error)

    def to_json(self):
        return {
            'LINE_NUM' : self.row,
            'ERROR_MSG': self.message
        }
